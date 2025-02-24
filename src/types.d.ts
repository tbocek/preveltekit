export type RouteParams = Record<string, string>;

export interface Route {
    path: string;
    component: Component;
}

export type Routes = Route[];
export type Component = any;