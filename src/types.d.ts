export type RouteParams = Record<string, string>;

export interface Route {
    path: string;
    component: Component;
    static: string;
}

export type Routes = Route[];
export type Component = any;