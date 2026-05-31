package validator

import (
	"errors"
	"regexp"
	"strconv"
	"time"
)

// 身份证号码验证规则
var (
	// 18位身份证号正则
	idCard18Regex = regexp.MustCompile(`^\d{17}[\dXx]$`)

	// 加权因子
	weightFactors = []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}

	// 校验码对应值
	checkCodes = []byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}
)

// ValidateIDCard 验证18位身份证号(ISO 7064 MOD 11-2)
func ValidateIDCard(idCard string) error {
	if len(idCard) != 18 {
		return errors.New("身份证号长度必须为18位")
	}

	if !idCard18Regex.MatchString(idCard) {
		return errors.New("身份证号格式错误")
	}

	// 验证出生日期
	birthDate := idCard[6:14]
	year, _ := strconv.Atoi(birthDate[0:4])
	month, _ := strconv.Atoi(birthDate[4:6])
	day, _ := strconv.Atoi(birthDate[6:8])

	if year < 1900 || year > time.Now().Year() {
		return errors.New("身份证号出生年份无效")
	}
	if month < 1 || month > 12 {
		return errors.New("身份证号出生月份无效")
	}
	if day < 1 || day > 31 {
		return errors.New("身份证号出生日期无效")
	}

	// 验证日期是否有效
	_, err := time.Parse("20060102", birthDate)
	if err != nil {
		return errors.New("身份证号出生日期无效")
	}

	// 计算校验码
	sum := 0
	for i := 0; i < 17; i++ {
		num, _ := strconv.Atoi(string(idCard[i]))
		sum += num * weightFactors[i]
	}
	expectedCheckCode := checkCodes[sum%11]

	// 验证校验码
	actualCheckCode := idCard[17]
	if actualCheckCode == 'x' {
		actualCheckCode = 'X'
	}

	if actualCheckCode != expectedCheckCode {
		return errors.New("身份证号校验码错误")
	}

	return nil
}

// ExtractIDCardInfo 从身份证号提取信息
func ExtractIDCardInfo(idCard string) (province string, birthday string, gender int, err error) {
	if err = ValidateIDCard(idCard); err != nil {
		return "", "", 0, err
	}

	// 提取省份代码
	provinceCode := idCard[0:2]
	province = getProvinceName(provinceCode)

	// 提取出生日期
	birthday = idCard[6:14]

	// 提取性别(倒数第二位奇数为男，偶数为女)
	genderCode, _ := strconv.Atoi(string(idCard[16]))
	if genderCode%2 == 1 {
		gender = 1 // 男
	} else {
		gender = 2 // 女
	}

	return province, birthday, gender, nil
}

// getProvinceName 根据省份代码获取省份名称
func getProvinceName(code string) string {
	provinces := map[string]string{
		"11": "北京", "12": "天津", "13": "河北", "14": "山西", "15": "内蒙古",
		"21": "辽宁", "22": "吉林", "23": "黑龙江",
		"31": "上海", "32": "江苏", "33": "浙江", "34": "安徽", "35": "福建", "36": "江西", "37": "山东",
		"41": "河南", "42": "湖北", "43": "湖南", "44": "广东", "45": "广西", "46": "海南",
		"50": "重庆", "51": "四川", "52": "贵州", "53": "云南", "54": "西藏",
		"61": "陕西", "62": "甘肃", "63": "青海", "64": "宁夏", "65": "新疆",
		"71": "台湾", "81": "香港", "82": "澳门", "91": "国外",
	}
	if name, ok := provinces[code]; ok {
		return name
	}
	return "未知"
}
