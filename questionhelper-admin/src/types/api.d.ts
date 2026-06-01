/** API 响应结构 */
export interface ApiResponse<T = any> {
  code: string;
  msg: string;
  data: T;
}

/** 分页响应 */
export interface PageResponse<T = any> {
  list: T[];
  total: number;
  page: number;
  page_size: number;
}

/** 分页参数 */
export interface PageParams {
  page: number;
  page_size: number;
}

/** 菜单项 */
export interface MenuItem {
  id: number;
  parent_id: number | null;
  name: string;
  title: string;
  type: 1 | 2 | 3; // 1=目录 2=菜单 3=按钮
  path: string;
  component: string;
  redirect: string;
  icon: string;
  hidden: boolean;
  permission: string;
  sort: number;
  status: number;
  children?: MenuItem[];
}

/** 路由项（前端格式） */
export interface RouteItem {
  name: string;
  path: string;
  component: string;
  redirect: string;
  meta: {
    title: string;
    icon: string;
    hidden: boolean;
    keepAlive: boolean;
  };
  children: RouteItem[];
}
