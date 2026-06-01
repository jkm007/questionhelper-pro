package question

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"questionhelper-server/internal/model"
	questionRepo "questionhelper-server/internal/repository/question"
	"questionhelper-server/pkg/logger"
)

const snapshotChangeLogPrefix = "[快照]"

// SaveSnapshot saves the current question state as a snapshot version.
// The snapshot is stored as a QuestionVersion record with a "[快照]" prefix in the ChangeLog.
func SaveSnapshot(questionID uint) error {
	question, err := questionRepo.FindByID(questionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("题目不存在")
		}
		return fmt.Errorf("查询题目失败: %w", err)
	}

	changelog := snapshotChangeLogPrefix + "手动保存快照"
	if err := CreateVersion(question, changelog, question.CreatorID); err != nil {
		return fmt.Errorf("保存快照失败: %w", err)
	}

	logger.Infof("题目 %d 保存快照成功", questionID)
	return nil
}

// GetSnapshots lists all snapshot versions for a question.
// Snapshots are identified by the "[快照]" prefix in the ChangeLog field.
func GetSnapshots(questionID uint) ([]model.QuestionVersion, error) {
	versions, err := questionRepo.FindVersionsByQuestionID(questionID)
	if err != nil {
		return nil, fmt.Errorf("查询版本列表失败: %w", err)
	}

	snapshots := make([]model.QuestionVersion, 0, len(versions))
	for _, v := range versions {
		if strings.HasPrefix(v.ChangeLog, snapshotChangeLogPrefix) {
			snapshots = append(snapshots, v)
		}
	}

	return snapshots, nil
}

// RestoreSnapshot restores a question to the state captured in a snapshot.
// It validates ownership, then delegates to RollbackVersion which handles
// saving the current state and applying the snapshot content.
func RestoreSnapshot(questionID, versionID uint) error {
	version, err := questionRepo.FindVersionByID(versionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("快照不存在")
		}
		return fmt.Errorf("查询快照失败: %w", err)
	}

	if version.QuestionID != questionID {
		return errors.New("快照不属于该题目")
	}

	// Delegate to RollbackVersion to avoid duplicating restore logic.
	// userID=0 indicates a system/snapshot restore.
	return RollbackVersion(int(questionID), version.Version, 0)
}
