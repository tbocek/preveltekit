"use strict";
(self["webpackChunkpreveltekit_example"] = self["webpackChunkpreveltekit_example"] || []).push([["298"], {
966: (function (__unused_webpack_module, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  A: () => (Router)
});
/* ESM import */var svelte_internal_disclose_version__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(999);
/* ESM import */var svelte_internal_client__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(750);
/* ESM import */var svelte__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(732);



var root_3 = svelte_internal_client__WEBPACK_IMPORTED_MODULE_1__/* .from_html */.vUu(`<h1> </h1>`);
function Router($$anchor, $$props) {
    svelte_internal_client__WEBPACK_IMPORTED_MODULE_1__/* .push */.VCO($$props, true);
    const renderer = ($$anchor)=>{
        var fragment = svelte_internal_client__WEBPACK_IMPORTED_MODULE_1__/* .comment */.Imx();
        var node = svelte_internal_client__WEBPACK_IMPORTED_MODULE_1__/* .first_child */.esp(fragment);
        {
            var consequent = ($$anchor)=>{
                var fragment_1 = svelte_internal_client__WEBPACK_IMPORTED_MODULE_1__/* .comment */.Imx();
                var node_1 = svelte_internal_client__WEBPACK_IMPORTED_MODULE_1__/* .first_child */.esp(fragment_1);
                svelte_internal_client__WEBPACK_IMPORTED_MODULE_1__/* .snippet */.UAl(node_1, ()=>svelte_internal_client__WEBPACK_IMPORTED_MODULE_1__/* .get */.JtY(activeComponent), ()=>({
                        params: svelte_internal_client__WEBPACK_IMPORTED_MODULE_1__/* .get */.JtY(routeParams)
                    }));
                svelte_internal_client__WEBPACK_IMPORTED_MODULE_1__/* .append */.BCw($$anchor, fragment_1);
            };
            var alternate = ($$anchor)=>{
                var h1 = root_3();
                var text = svelte_internal_client__WEBPACK_IMPORTED_MODULE_1__/* .child */.jfp(h1);
                svelte_internal_client__WEBPACK_IMPORTED_MODULE_1__/* .reset */.cLc(h1);
                svelte_internal_client__WEBPACK_IMPORTED_MODULE_1__/* .template_effect */.vNg(()=>svelte_internal_client__WEBPACK_IMPORTED_MODULE_1__/* .set_text */.jax(text, `404 - Page Not Found for [${svelte_internal_client__WEBPACK_IMPORTED_MODULE_1__/* .get */.JtY(currentRoute) ?? ''}]`));
                svelte_internal_client__WEBPACK_IMPORTED_MODULE_1__/* .append */.BCw($$anchor, h1);
            };
            svelte_internal_client__WEBPACK_IMPORTED_MODULE_1__["if"](node, ($$render)=>{
                if (svelte_internal_client__WEBPACK_IMPORTED_MODULE_1__/* .get */.JtY(activeComponent)) $$render(consequent);
                else $$render(alternate, false);
            });
        }
        svelte_internal_client__WEBPACK_IMPORTED_MODULE_1__/* .append */.BCw($$anchor, fragment);
    };
    function navigate(path) {
        history.pushState(null, "", path);
        window.dispatchEvent(new CustomEvent('svelteNavigate', {
            detail: {
                path
            }
        }));
    }
    // Props for the component
    const props = svelte_internal_client__WEBPACK_IMPORTED_MODULE_1__/* .rest_props */.iRd($$props, [
        '$$slots',
        '$$events',
        '$$legacy'
    ]);
    // Create a store for the current route
    let currentRoute = svelte_internal_client__WEBPACK_IMPORTED_MODULE_1__/* .state */.wk1('/');
    // Extract routes and notFound from props with defaults
    let routes = svelte_internal_client__WEBPACK_IMPORTED_MODULE_1__/* .derived */.unG(()=>$$props.routes || []);
    // Handle browser back/forward navigation
    const handlePopState = ()=>{
        svelte_internal_client__WEBPACK_IMPORTED_MODULE_1__/* .set */.hZp(currentRoute, window.location.pathname, true);
    };
    // Define event handler function for custom navigation events
    const handleNavigateEvent = (e)=>{
        const customEvent = e;
        svelte_internal_client__WEBPACK_IMPORTED_MODULE_1__/* .set */.hZp(currentRoute, customEvent.detail.path, true);
    };
    (0,svelte__WEBPACK_IMPORTED_MODULE_2__/* .onMount */.Rc)(()=>{
        var _window;
        // Set initial route
        svelte_internal_client__WEBPACK_IMPORTED_MODULE_1__/* .set */.hZp(currentRoute, window.location.pathname, true);
        // Add event listener for back/forward navigation
        window.addEventListener('popstate', handlePopState);
        window.addEventListener('svelteNavigate', handleNavigateEvent);
        // Expose routes directly to SSR if running in JSDOM
        if ((_window = window) === null || _window === void 0 ? void 0 : _window.JSDOM) {
            // Get the current value of the routes
            const currentRoutes = [
                ...svelte_internal_client__WEBPACK_IMPORTED_MODULE_1__/* .get */.JtY(routes)
            ];
            window.__svelteRoutes = currentRoutes;
        }
    });
    (0,svelte__WEBPACK_IMPORTED_MODULE_2__/* .onDestroy */.sA)(()=>{
        window.removeEventListener('popstate', handlePopState);
        window.removeEventListener('svelteNavigate', handleNavigateEvent);
    });
    // Helper function to find matching route and extract params
    function findMatchingRoute(path) {
        // Normalize the input path (remove trailing slash except for root path)
        const normalizedPath = path === '/' ? '/' : path.endsWith('/') ? path.slice(0, -1) : path;
        // For storing the best match
        let bestMatch = null;
        // Check all routes for matches
        for (const route of svelte_internal_client__WEBPACK_IMPORTED_MODULE_1__/* .get */.JtY(routes)){
            const routePath = route.path;
            let isMatch = false;
            let params = {};
            let specificity = 0;
            // Normalize the route path (remove trailing slash except for root path)
            const normalizedRoutePath = routePath === '/' ? '/' : routePath.endsWith('/') ? routePath.slice(0, -1) : routePath;
            // CASE 1: Handle root path special case
            if (normalizedRoutePath === '/') {
                isMatch = normalizedPath === '/';
                specificity = 100; // Highest specificity for root path
            } else if (normalizedRoutePath === '**/' || normalizedRoutePath === '**/') {
                isMatch = normalizedPath === '/';
                specificity = 1; // Lowest specificity
            } else if (normalizedRoutePath === '*/' || normalizedRoutePath === '*') {
                // Match root path
                if (normalizedPath === '/') {
                    isMatch = true;
                    specificity = 1;
                } else {
                    const pathSegments = normalizedPath.split('/').filter(Boolean);
                    if (pathSegments.length === 1) {
                        isMatch = true;
                        specificity = 1;
                    }
                }
            } else if (normalizedRoutePath.startsWith('*/')) {
                // Get the suffix after */
                const suffix = normalizedRoutePath.slice(2);
                // Check if path exactly matches suffix
                if (normalizedPath === '/' + suffix) {
                    isMatch = true;
                    specificity = 2;
                } else {
                    // Build a regex that matches /{segment}/{suffix} exactly
                    const pattern = new RegExp(`^\\/([^\\/]+)\\/${suffix}$`);
                    isMatch = pattern.test(normalizedPath);
                    specificity = 2;
                }
            } else if (normalizedRoutePath.startsWith('*')) {
                const suffix = normalizedRoutePath.slice(1);
                // Simple wildcard matching for other patterns
                if (normalizedPath === suffix || suffix.startsWith('/') && normalizedPath.endsWith(suffix)) {
                    isMatch = true;
                    specificity = 2; // Low specificity for wildcard routes
                }
            } else {
                // Split paths into segments for comparison
                const routeSegments = normalizedRoutePath.split('/').filter(Boolean);
                const pathSegments = normalizedPath.split('/').filter(Boolean);
                // For standard routes, segment count must match
                if (routeSegments.length === pathSegments.length) {
                    isMatch = true;
                    specificity = 0;
                    // Compare each segment
                    for(let i = 0; i < routeSegments.length; i++){
                        const routeSegment = routeSegments[i];
                        const pathSegment = pathSegments[i];
                        // Handle parameter segments
                        if (routeSegment.startsWith(':')) {
                            const paramName = routeSegment.slice(1);
                            params[paramName] = pathSegment;
                            specificity += 5;
                        } else if (routeSegment === pathSegment) {
                            specificity += 10;
                        } else {
                            isMatch = false;
                            break;
                        }
                    }
                }
            }
            // Update best match if this route matches and is more specific
            if (isMatch && (!bestMatch || specificity > bestMatch.specificity)) {
                bestMatch = {
                    component: route.component,
                    params,
                    specificity
                };
            }
        }
        // Return the best match or null component
        return bestMatch ? {
            component: bestMatch.component,
            params: bestMatch.params
        } : {
            component: null,
            params: {}
        };
    }
    // Current matched route and params
    let matchedRoute = svelte_internal_client__WEBPACK_IMPORTED_MODULE_1__/* .derived */.unG(()=>findMatchingRoute(svelte_internal_client__WEBPACK_IMPORTED_MODULE_1__/* .get */.JtY(currentRoute)));
    let activeComponent = svelte_internal_client__WEBPACK_IMPORTED_MODULE_1__/* .derived */.unG(()=>svelte_internal_client__WEBPACK_IMPORTED_MODULE_1__/* .get */.JtY(matchedRoute).component);
    let routeParams = svelte_internal_client__WEBPACK_IMPORTED_MODULE_1__/* .derived */.unG(()=>svelte_internal_client__WEBPACK_IMPORTED_MODULE_1__/* .get */.JtY(matchedRoute).params);
    var $$exports = {
        navigate
    };
    renderer($$anchor);
    return svelte_internal_client__WEBPACK_IMPORTED_MODULE_1__/* .pop */.uYY($$exports);
}


}),
178: (function (__unused_webpack_module, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  Ax: () => (TEMPLATE_FRAGMENT),
  CD: () => (HYDRATION_START),
  Lc: () => (HYDRATION_END),
  UP: () => (UNINITIALIZED),
  Uh: () => (FILENAME),
  iX: () => (TEMPLATE_USE_IMPORT_NODE),
  kD: () => (HYDRATION_ERROR),
  qn: () => (HYDRATION_START_ELSE)
});
const EACH_ITEM_REACTIVE = 1;
const EACH_INDEX_REACTIVE = (/* unused pure expression or super */ null && (1 << 1));
/** See EachBlock interface metadata.is_controlled for an explanation what this is */ const EACH_IS_CONTROLLED = (/* unused pure expression or super */ null && (1 << 2));
const EACH_IS_ANIMATED = (/* unused pure expression or super */ null && (1 << 3));
const EACH_ITEM_IMMUTABLE = (/* unused pure expression or super */ null && (1 << 4));
const PROPS_IS_IMMUTABLE = 1;
const PROPS_IS_RUNES = (/* unused pure expression or super */ null && (1 << 1));
const PROPS_IS_UPDATED = (/* unused pure expression or super */ null && (1 << 2));
const PROPS_IS_BINDABLE = (/* unused pure expression or super */ null && (1 << 3));
const PROPS_IS_LAZY_INITIAL = (/* unused pure expression or super */ null && (1 << 4));
const TRANSITION_IN = 1;
const TRANSITION_OUT = (/* unused pure expression or super */ null && (1 << 1));
const TRANSITION_GLOBAL = (/* unused pure expression or super */ null && (1 << 2));
const TEMPLATE_FRAGMENT = 1;
const TEMPLATE_USE_IMPORT_NODE = 1 << 1;
const TEMPLATE_USE_SVG = (/* unused pure expression or super */ null && (1 << 2));
const TEMPLATE_USE_MATHML = (/* unused pure expression or super */ null && (1 << 3));
const HYDRATION_START = '[';
/** used to indicate that an `{:else}...` block was rendered */ const HYDRATION_START_ELSE = '[!';
const HYDRATION_END = ']';
const HYDRATION_ERROR = {};
const ELEMENT_IS_NAMESPACED = 1;
const ELEMENT_PRESERVE_ATTRIBUTE_CASE = (/* unused pure expression or super */ null && (1 << 1));
const ELEMENT_IS_INPUT = (/* unused pure expression or super */ null && (1 << 2));
const UNINITIALIZED = Symbol();
// Dev-time component properties
const FILENAME = Symbol('filename');
const HMR = Symbol('hmr');
const NAMESPACE_HTML = 'http://www.w3.org/1999/xhtml';
const NAMESPACE_SVG = 'http://www.w3.org/2000/svg';
const NAMESPACE_MATHML = 'http://www.w3.org/1998/Math/MathML';
// we use a list of ignorable runtime warnings because not every runtime warning
// can be ignored and we want to keep the validation for svelte-ignore in place
const IGNORABLE_RUNTIME_WARNINGS = /** @type {const} */ (/* unused pure expression or super */ null && ([
    'await_waterfall',
    'await_reactivity_loss',
    'state_snapshot_uncloneable',
    'binding_property_non_reactive',
    'hydration_attribute_changed',
    'hydration_html_changed',
    'ownership_invalid_binding',
    'ownership_invalid_mutation'
]));
/**
 * Whitespace inside one of these elements will not result in
 * a whitespace node being created in any circumstances. (This
 * list is almost certainly very incomplete)
 * TODO this is currently unused
 */ const ELEMENTS_WITHOUT_TEXT = (/* unused pure expression or super */ null && ([
    'audio',
    'datalist',
    'dl',
    'optgroup',
    'select',
    'video'
]));
const ATTACHMENT_KEY = '@attach';


}),
732: (function (__unused_webpack_module, __webpack_exports__, __webpack_require__) {

// EXPORTS
__webpack_require__.d(__webpack_exports__, {
  sA: () => (/* binding */ onDestroy),
  Qv: () => (/* reexport */ render/* .hydrate */.Qv),
  Rc: () => (/* binding */ onMount)
});

// UNUSED EXPORTS: getAllContexts, hasContext, beforeUpdate, createRawSnippet, untrack, setContext, flushSync, afterUpdate, getContext, mount, getAbortSignal, tick, createEventDispatcher, unmount, settled

// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/runtime.js
var runtime = __webpack_require__(513);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/shared/utils.js
var utils = __webpack_require__(986);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/index.js + 50 modules
var client = __webpack_require__(750);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/errors.js
var errors = __webpack_require__(626);
// EXTERNAL MODULE: ./node_modules/.pnpm/esm-env@1.2.2/node_modules/esm-env/false.js
var esm_env_false = __webpack_require__(832);
;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/shared/errors.js
/* This file is generated by scripts/process-messages/index.js. Do not edit! */ 
/**
 * Cannot use `{@render children(...)}` if the parent component uses `let:` directives. Consider using a named snippet instead
 * @returns {never}
 */ function invalid_default_snippet() {
    if (DEV) {
        const error = new Error(`invalid_default_snippet\nCannot use \`{@render children(...)}\` if the parent component uses \`let:\` directives. Consider using a named snippet instead\nhttps://svelte.dev/e/invalid_default_snippet`);
        error.name = 'Svelte error';
        throw error;
    } else {
        throw new Error(`https://svelte.dev/e/invalid_default_snippet`);
    }
}
/**
 * A snippet function was passed invalid arguments. Snippets should only be instantiated via `{@render ...}`
 * @returns {never}
 */ function invalid_snippet_arguments() {
    if (DEV) {
        const error = new Error(`invalid_snippet_arguments\nA snippet function was passed invalid arguments. Snippets should only be instantiated via \`{@render ...}\`\nhttps://svelte.dev/e/invalid_snippet_arguments`);
        error.name = 'Svelte error';
        throw error;
    } else {
        throw new Error(`https://svelte.dev/e/invalid_snippet_arguments`);
    }
}
/**
 * `%name%(...)` can only be used during component initialisation
 * @param {string} name
 * @returns {never}
 */ function lifecycle_outside_component(name) {
    if (esm_env_false/* ["default"] */.A) {
        const error = new Error(`lifecycle_outside_component\n\`${name}(...)\` can only be used during component initialisation\nhttps://svelte.dev/e/lifecycle_outside_component`);
        error.name = 'Svelte error';
        throw error;
    } else {
        throw new Error(`https://svelte.dev/e/lifecycle_outside_component`);
    }
}
/**
 * Attempted to render a snippet without a `{@render}` block. This would cause the snippet code to be stringified instead of its content being rendered to the DOM. To fix this, change `{snippet}` to `{@render snippet()}`.
 * @returns {never}
 */ function snippet_without_render_tag() {
    if (DEV) {
        const error = new Error(`snippet_without_render_tag\nAttempted to render a snippet without a \`{@render}\` block. This would cause the snippet code to be stringified instead of its content being rendered to the DOM. To fix this, change \`{snippet}\` to \`{@render snippet()}\`.\nhttps://svelte.dev/e/snippet_without_render_tag`);
        error.name = 'Svelte error';
        throw error;
    } else {
        throw new Error(`https://svelte.dev/e/snippet_without_render_tag`);
    }
}
/**
 * `%name%` is not a store with a `subscribe` method
 * @param {string} name
 * @returns {never}
 */ function store_invalid_shape(name) {
    if (DEV) {
        const error = new Error(`store_invalid_shape\n\`${name}\` is not a store with a \`subscribe\` method\nhttps://svelte.dev/e/store_invalid_shape`);
        error.name = 'Svelte error';
        throw error;
    } else {
        throw new Error(`https://svelte.dev/e/store_invalid_shape`);
    }
}
/**
 * The `this` prop on `<svelte:element>` must be a string, if defined
 * @returns {never}
 */ function svelte_element_invalid_this_value() {
    if (DEV) {
        const error = new Error(`svelte_element_invalid_this_value\nThe \`this\` prop on \`<svelte:element>\` must be a string, if defined\nhttps://svelte.dev/e/svelte_element_invalid_this_value`);
        error.name = 'Svelte error';
        throw error;
    } else {
        throw new Error(`https://svelte.dev/e/svelte_element_invalid_this_value`);
    }
}

// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/flags/index.js
var flags = __webpack_require__(817);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/context.js
var client_context = __webpack_require__(754);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/reactivity/batch.js
var batch = __webpack_require__(410);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/render.js
var render = __webpack_require__(485);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/blocks/snippet.js
var snippet = __webpack_require__(768);
;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/index-client.js
/** @import { ComponentContext, ComponentContextLegacy } from '#client' */ /** @import { EventDispatcher } from './index.js' */ /** @import { NotFunction } from './internal/types.js' */ 






if (esm_env_false/* ["default"] */.A) {
    /**
	 * @param {string} rune
	 */ function throw_rune_error(rune) {
        if (!(rune in globalThis)) {
            // TODO if people start adjusting the "this can contain runes" config through v-p-s more, adjust this message
            /** @type {any} */ let value; // let's hope noone modifies this global, but belts and braces
            Object.defineProperty(globalThis, rune, {
                configurable: true,
                // eslint-disable-next-line getter-return
                get: ()=>{
                    if (value !== undefined) {
                        return value;
                    }
                    errors/* .rune_outside_svelte */.xU(rune);
                },
                set: (v)=>{
                    value = v;
                }
            });
        }
    }
    throw_rune_error('$state');
    throw_rune_error('$effect');
    throw_rune_error('$derived');
    throw_rune_error('$inspect');
    throw_rune_error('$props');
    throw_rune_error('$bindable');
}
/**
 * Returns an [`AbortSignal`](https://developer.mozilla.org/en-US/docs/Web/API/AbortSignal) that aborts when the current [derived](https://svelte.dev/docs/svelte/$derived) or [effect](https://svelte.dev/docs/svelte/$effect) re-runs or is destroyed.
 *
 * Must be called while a derived or effect is running.
 *
 * ```svelte
 * <script>
 * 	import { getAbortSignal } from 'svelte';
 *
 * 	let { id } = $props();
 *
 * 	async function getData(id) {
 * 		const response = await fetch(`/items/${id}`, {
 * 			signal: getAbortSignal()
 * 		});
 *
 * 		return await response.json();
 * 	}
 *
 * 	const data = $derived(await getData(id));
 * </script>
 * ```
 */ function getAbortSignal() {
    var _active_reaction;
    if (active_reaction === null) {
        e.get_abort_signal_outside_reaction();
    }
    return ((_active_reaction = active_reaction).ac ?? (_active_reaction.ac = new AbortController())).signal;
}
/**
 * `onMount`, like [`$effect`](https://svelte.dev/docs/svelte/$effect), schedules a function to run as soon as the component has been mounted to the DOM.
 * Unlike `$effect`, the provided function only runs once.
 *
 * It must be called during the component's initialisation (but doesn't need to live _inside_ the component;
 * it can be called from an external module). If a function is returned _synchronously_ from `onMount`,
 * it will be called when the component is unmounted.
 *
 * `onMount` functions do not run during [server-side rendering](https://svelte.dev/docs/svelte/svelte-server#render).
 *
 * @template T
 * @param {() => NotFunction<T> | Promise<NotFunction<T>> | (() => any)} fn
 * @returns {void}
 */ function onMount(fn) {
    if (client_context/* .component_context */.UL === null) {
        lifecycle_outside_component('onMount');
    }
    if (flags/* .legacy_mode_flag */.LM && client_context/* .component_context.l */.UL.l !== null) {
        init_update_callbacks(client_context/* .component_context */.UL).m.push(fn);
    } else {
        (0,client/* .user_effect */.MWq)(()=>{
            const cleanup = (0,runtime/* .untrack */.vz)(fn);
            if (typeof cleanup === 'function') return /** @type {() => void} */ cleanup;
        });
    }
}
/**
 * Schedules a callback to run immediately before the component is unmounted.
 *
 * Out of `onMount`, `beforeUpdate`, `afterUpdate` and `onDestroy`, this is the
 * only one that runs inside a server-side component.
 *
 * @param {() => any} fn
 * @returns {void}
 */ function onDestroy(fn) {
    if (client_context/* .component_context */.UL === null) {
        lifecycle_outside_component('onDestroy');
    }
    onMount(()=>()=>(0,runtime/* .untrack */.vz)(fn));
}
/**
 * @template [T=any]
 * @param {string} type
 * @param {T} [detail]
 * @param {any}params_0
 * @returns {CustomEvent<T>}
 */ function create_custom_event(type, detail) {
    let { bubbles = false, cancelable = false } = arguments.length > 2 && arguments[2] !== void 0 ? arguments[2] : {};
    return new CustomEvent(type, {
        detail,
        bubbles,
        cancelable
    });
}
/**
 * Creates an event dispatcher that can be used to dispatch [component events](https://svelte.dev/docs/svelte/legacy-on#Component-events).
 * Event dispatchers are functions that can take two arguments: `name` and `detail`.
 *
 * Component events created with `createEventDispatcher` create a
 * [CustomEvent](https://developer.mozilla.org/en-US/docs/Web/API/CustomEvent).
 * These events do not [bubble](https://developer.mozilla.org/en-US/docs/Learn/JavaScript/Building_blocks/Events#Event_bubbling_and_capture).
 * The `detail` argument corresponds to the [CustomEvent.detail](https://developer.mozilla.org/en-US/docs/Web/API/CustomEvent/detail)
 * property and can contain any type of data.
 *
 * The event dispatcher can be typed to narrow the allowed event names and the type of the `detail` argument:
 * ```ts
 * const dispatch = createEventDispatcher<{
 *  loaded: null; // does not take a detail argument
 *  change: string; // takes a detail argument of type string, which is required
 *  optional: number | null; // takes an optional detail argument of type number
 * }>();
 * ```
 *
 * @deprecated Use callback props and/or the `$host()` rune instead — see [migration guide](https://svelte.dev/docs/svelte/v5-migration-guide#Event-changes-Component-events)
 * @template {Record<string, any>} [EventMap = any]
 * @returns {EventDispatcher<EventMap>}
 */ function createEventDispatcher() {
    const active_component_context = component_context;
    if (active_component_context === null) {
        e.lifecycle_outside_component('createEventDispatcher');
    }
    /**
	 * @param [detail]
	 * @param [options]
	 */ return (type, detail, options)=>{
        var /** @type {Record<string, Function | Function[]>} */ _active_component_context_s_$$events;
        const events = (_active_component_context_s_$$events = active_component_context.s.$$events) === null || _active_component_context_s_$$events === void 0 ? void 0 : _active_component_context_s_$$events[/** @type {string} */ type];
        if (events) {
            const callbacks = is_array(events) ? events.slice() : [
                events
            ];
            // TODO are there situations where events could be dispatched
            // in a server (non-DOM) environment?
            const event = create_custom_event(/** @type {string} */ type, detail, options);
            for (const fn of callbacks){
                fn.call(active_component_context.x, event);
            }
            return !event.defaultPrevented;
        }
        return true;
    };
}
// TODO mark beforeUpdate and afterUpdate as deprecated in Svelte 6
/**
 * Schedules a callback to run immediately before the component is updated after any state change.
 *
 * The first time the callback runs will be before the initial `onMount`.
 *
 * In runes mode use `$effect.pre` instead.
 *
 * @deprecated Use [`$effect.pre`](https://svelte.dev/docs/svelte/$effect#$effect.pre) instead
 * @param {() => void} fn
 * @returns {void}
 */ function beforeUpdate(fn) {
    if (component_context === null) {
        e.lifecycle_outside_component('beforeUpdate');
    }
    if (component_context.l === null) {
        e.lifecycle_legacy_only('beforeUpdate');
    }
    init_update_callbacks(component_context).b.push(fn);
}
/**
 * Schedules a callback to run immediately after the component has been updated.
 *
 * The first time the callback runs will be after the initial `onMount`.
 *
 * In runes mode use `$effect` instead.
 *
 * @deprecated Use [`$effect`](https://svelte.dev/docs/svelte/$effect) instead
 * @param {() => void} fn
 * @returns {void}
 */ function afterUpdate(fn) {
    if (component_context === null) {
        e.lifecycle_outside_component('afterUpdate');
    }
    if (component_context.l === null) {
        e.lifecycle_legacy_only('afterUpdate');
    }
    init_update_callbacks(component_context).a.push(fn);
}
/**
 * Legacy-mode: Init callbacks object for onMount/beforeUpdate/afterUpdate
 * @param {ComponentContext} context
 */ function init_update_callbacks(context) {
    var _l;
    var l = /** @type {ComponentContextLegacy} */ context.l;
    return (_l = l).u ?? (_l.u = {
        a: [],
        b: [],
        m: []
    });
}







}),
924: (function (__unused_webpack_module, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  $q: () => (INERT),
  EY: () => (REACTION_IS_UPDATING),
  FV: () => (ROOT_EFFECT),
  In: () => (STALE_REACTION),
  L2: () => (UNOWNED),
  Nd: () => (TEXT_NODE),
  PL: () => (HEAD_EFFECT),
  Qf: () => (PROXY_PATH_SYMBOL),
  T1: () => (INSPECT_EFFECT),
  V$: () => (EFFECT_PRESERVED),
  VD: () => (ASYNC),
  Wr: () => (USER_EFFECT),
  Zr: () => (BRANCH_EFFECT),
  Zv: () => (RENDER_EFFECT),
  _N: () => (DISCONNECTED),
  ac: () => (EFFECT),
  bp: () => (BOUNDARY_EFFECT),
  dH: () => (ERROR_VALUE),
  dz: () => (COMMENT_NODE),
  ig: () => (MAYBE_DIRTY),
  jm: () => (DIRTY),
  kc: () => (BLOCK_EFFECT),
  l3: () => (LEGACY_PROPS),
  lQ: () => (EFFECT_TRANSPARENT),
  mj: () => (DERIVED),
  o5: () => (DESTROYED),
  w_: () => (CLEAN),
  wi: () => (EFFECT_RAN),
  x3: () => (STATE_SYMBOL)
});
/* ESM import */var _swc_helpers_define_property__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(925);

const DERIVED = 1 << 1;
const EFFECT = 1 << 2;
const RENDER_EFFECT = 1 << 3;
const BLOCK_EFFECT = 1 << 4;
const BRANCH_EFFECT = 1 << 5;
const ROOT_EFFECT = 1 << 6;
const BOUNDARY_EFFECT = 1 << 7;
const UNOWNED = 1 << 8;
const DISCONNECTED = 1 << 9;
const CLEAN = 1 << 10;
const DIRTY = 1 << 11;
const MAYBE_DIRTY = 1 << 12;
const INERT = 1 << 13;
const DESTROYED = 1 << 14;
const EFFECT_RAN = 1 << 15;
/** 'Transparent' effects do not create a transition boundary */ const EFFECT_TRANSPARENT = 1 << 16;
const INSPECT_EFFECT = 1 << 17;
const HEAD_EFFECT = 1 << 18;
const EFFECT_PRESERVED = 1 << 19;
const USER_EFFECT = 1 << 20;
// Flags used for async
const REACTION_IS_UPDATING = 1 << 21;
const ASYNC = 1 << 22;
const ERROR_VALUE = 1 << 23;
const STATE_SYMBOL = Symbol('$state');
const LEGACY_PROPS = Symbol('legacy props');
const LOADING_ATTR_SYMBOL = Symbol('');
const PROXY_PATH_SYMBOL = Symbol('proxy path');
/** allow users to ignore aborted signal errors if `reason.name === 'StaleReactionError` */ const STALE_REACTION = new class StaleReactionError extends Error {
    constructor(...args){
        super(...args), (0,_swc_helpers_define_property__WEBPACK_IMPORTED_MODULE_0__._)(this, "name", 'StaleReactionError'), (0,_swc_helpers_define_property__WEBPACK_IMPORTED_MODULE_0__._)(this, "message", 'The reaction that called `getAbortSignal()` was re-run or destroyed');
    }
}();
const ELEMENT_NODE = 1;
const TEXT_NODE = 3;
const COMMENT_NODE = 8;
const DOCUMENT_FRAGMENT_NODE = 11;


}),
754: (function (__unused_webpack_module, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  DE: () => (dev_current_component_function),
  De: () => (set_component_context),
  Mo: () => (set_dev_current_component_function),
  O2: () => (set_dev_stack),
  UL: () => (component_context),
  VC: () => (push),
  hH: () => (is_runes),
  lv: () => (dev_stack),
  uY: () => (pop)
});
/* ESM import */var esm_env__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(832);
/* ESM import */var _runtime_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(513);
/* ESM import */var _reactivity_effects_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(480);
/* ESM import */var _flags_index_js__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(817);
/* ESM import */var _constants_js__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(178);
/* ESM import */var _constants_js__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(924);
/** @import { ComponentContext, DevStackEntry, Effect } from '#client' */ 






/** @type {ComponentContext | null} */ let component_context = null;
/** @param {ComponentContext | null} context */ function set_component_context(context) {
    component_context = context;
}
/** @type {DevStackEntry | null} */ let dev_stack = null;
/** @param {DevStackEntry | null} stack */ function set_dev_stack(stack) {
    dev_stack = stack;
}
/**
 * Execute a callback with a new dev stack entry
 * @param {() => any} callback - Function to execute
 * @param {DevStackEntry['type']} type - Type of block/component
 * @param {any} component - Component function
 * @param {number} line - Line number
 * @param {number} column - Column number
 * @param {Record<string, any>} [additional] - Any additional properties to add to the dev stack entry
 * @returns {any}
 */ function add_svelte_meta(callback, type, component, line, column, additional) {
    const parent = dev_stack;
    dev_stack = {
        type,
        file: component[FILENAME],
        line,
        column,
        parent,
        ...additional
    };
    try {
        return callback();
    } finally{
        dev_stack = parent;
    }
}
/**
 * The current component function. Different from current component context:
 * ```html
 * <!-- App.svelte -->
 * <Foo>
 *   <Bar /> <!-- context == Foo.svelte, function == App.svelte -->
 * </Foo>
 * ```
 * @type {ComponentContext['function']}
 */ let dev_current_component_function = null;
/** @param {ComponentContext['function']} fn */ function set_dev_current_component_function(fn) {
    dev_current_component_function = fn;
}
/**
 * Retrieves the context that belongs to the closest parent component with the specified `key`.
 * Must be called during component initialisation.
 *
 * @template T
 * @param {any} key
 * @returns {T}
 */ function getContext(key) {
    const context_map = get_or_init_context_map('getContext');
    const result = /** @type {T} */ context_map.get(key);
    return result;
}
/**
 * Associates an arbitrary `context` object with the current component and the specified `key`
 * and returns that object. The context is then available to children of the component
 * (including slotted content) with `getContext`.
 *
 * Like lifecycle functions, this must be called during component initialisation.
 *
 * @template T
 * @param {any} key
 * @param {T} context
 * @returns {T}
 */ function setContext(key, context) {
    const context_map = get_or_init_context_map('setContext');
    if (async_mode_flag) {
        var flags = /** @type {Effect} */ active_effect.f;
        var valid = !active_reaction && (flags & BRANCH_EFFECT) !== 0 && (flags & EFFECT_RAN) === 0;
        if (!valid) {
            e.set_context_after_init();
        }
    }
    context_map.set(key, context);
    return context;
}
/**
 * Checks whether a given `key` has been set in the context of a parent component.
 * Must be called during component initialisation.
 *
 * @param {any} key
 * @returns {boolean}
 */ function hasContext(key) {
    const context_map = get_or_init_context_map('hasContext');
    return context_map.has(key);
}
/**
 * Retrieves the whole context map that belongs to the closest parent component.
 * Must be called during component initialisation. Useful, for example, if you
 * programmatically create a component and want to pass the existing context to it.
 *
 * @template {Map<any, any>} [T=Map<any, any>]
 * @returns {T}
 */ function getAllContexts() {
    const context_map = get_or_init_context_map('getAllContexts');
    return /** @type {T} */ context_map;
}
/**
 * @param {Record<string, unknown>} props
 * @param {any} runes
 * @param {Function} [fn]
 * @returns {void}
 */ function push(props) {
    let runes = arguments.length > 1 && arguments[1] !== void 0 ? arguments[1] : false, fn = arguments.length > 2 ? arguments[2] : void 0;
    component_context = {
        p: component_context,
        c: null,
        e: null,
        s: props,
        x: null,
        l: _flags_index_js__WEBPACK_IMPORTED_MODULE_4__/* .legacy_mode_flag */.LM && !runes ? {
            s: null,
            u: null,
            $: []
        } : null
    };
    if (esm_env__WEBPACK_IMPORTED_MODULE_5__/* ["default"] */.A) {
        // component function
        component_context.function = fn;
        dev_current_component_function = fn;
    }
}
/**
 * @template {Record<string, any>} T
 * @param {T} [component]
 * @returns {T}
 */ function pop(component) {
    var context = /** @type {ComponentContext} */ component_context;
    var effects = context.e;
    if (effects !== null) {
        context.e = null;
        for (var fn of effects){
            (0,_reactivity_effects_js__WEBPACK_IMPORTED_MODULE_1__/* .create_user_effect */.V1)(fn);
        }
    }
    if (component !== undefined) {
        context.x = component;
    }
    component_context = context.p;
    if (esm_env__WEBPACK_IMPORTED_MODULE_5__/* ["default"] */.A) {
        dev_current_component_function = (component_context === null || component_context === void 0 ? void 0 : component_context.function) ?? null;
    }
    return component ?? /** @type {T} */ {};
}
/** @returns {boolean} */ function is_runes() {
    return !_flags_index_js__WEBPACK_IMPORTED_MODULE_4__/* .legacy_mode_flag */.LM || component_context !== null && component_context.l === null;
}
/**
 * @param {string} name
 * @returns {Map<unknown, unknown>}
 */ function get_or_init_context_map(name) {
    var _component_context;
    if (component_context === null) {
        e.lifecycle_outside_component(name);
    }
    return (_component_context = component_context).c ?? (_component_context.c = new Map(get_parent_context(component_context) || undefined));
}
/**
 * @param {ComponentContext} component_context
 * @returns {Map<unknown, unknown> | null}
 */ function get_parent_context(component_context) {
    let parent = component_context.p;
    while(parent !== null){
        const context_map = parent.c;
        if (context_map !== null) {
            return context_map;
        }
        parent = parent.p;
    }
    return null;
}


}),
301: (function (__unused_webpack_module, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  Ej: () => (init_array_prototype_warnings)
});
/* ESM import */var _warnings_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(32);
/* ESM import */var _proxy_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(445);


function init_array_prototype_warnings() {
    const array_prototype = Array.prototype;
    // The REPL ends up here over and over, and this prevents it from adding more and more patches
    // of the same kind to the prototype, which would slow down everything over time.
    // @ts-expect-error
    const cleanup = Array.__svelte_cleanup;
    if (cleanup) {
        cleanup();
    }
    const { indexOf, lastIndexOf, includes } = array_prototype;
    array_prototype.indexOf = function(item, from_index) {
        const index = indexOf.call(this, item, from_index);
        if (index === -1) {
            for(let i = from_index ?? 0; i < this.length; i += 1){
                if ((0,_proxy_js__WEBPACK_IMPORTED_MODULE_0__/* .get_proxied_value */.N)(this[i]) === item) {
                    _warnings_js__WEBPACK_IMPORTED_MODULE_1__/* .state_proxy_equality_mismatch */.ns('array.indexOf(...)');
                    break;
                }
            }
        }
        return index;
    };
    array_prototype.lastIndexOf = function(item, from_index) {
        // we need to specify this.length - 1 because it's probably using something like
        // `arguments` inside so passing undefined is different from not passing anything
        const index = lastIndexOf.call(this, item, from_index ?? this.length - 1);
        if (index === -1) {
            for(let i = 0; i <= (from_index ?? this.length - 1); i += 1){
                if ((0,_proxy_js__WEBPACK_IMPORTED_MODULE_0__/* .get_proxied_value */.N)(this[i]) === item) {
                    _warnings_js__WEBPACK_IMPORTED_MODULE_1__/* .state_proxy_equality_mismatch */.ns('array.lastIndexOf(...)');
                    break;
                }
            }
        }
        return index;
    };
    array_prototype.includes = function(item, from_index) {
        const has = includes.call(this, item, from_index);
        if (!has) {
            for(let i = 0; i < this.length; i += 1){
                if ((0,_proxy_js__WEBPACK_IMPORTED_MODULE_0__/* .get_proxied_value */.N)(this[i]) === item) {
                    _warnings_js__WEBPACK_IMPORTED_MODULE_1__/* .state_proxy_equality_mismatch */.ns('array.includes(...)');
                    break;
                }
            }
        }
        return has;
    };
    // @ts-expect-error
    Array.__svelte_cleanup = ()=>{
        array_prototype.indexOf = indexOf;
        array_prototype.lastIndexOf = lastIndexOf;
        array_prototype.includes = includes;
    };
}
/**
 * @param {any} a
 * @param {any} b
 * @param {boolean} equal
 * @returns {boolean}
 */ function strict_equals(a, b) {
    let equal = arguments.length > 2 && arguments[2] !== void 0 ? arguments[2] : true;
    // try-catch needed because this tries to read properties of `a` and `b`,
    // which could be disallowed for example in a secure context
    try {
        if (a === b !== (get_proxied_value(a) === get_proxied_value(b))) {
            w.state_proxy_equality_mismatch(equal ? '===' : '!==');
        }
    } catch  {}
    return a === b === equal;
}
/**
 * @param {any} a
 * @param {any} b
 * @param {boolean} equal
 * @returns {boolean}
 */ function equals(a, b) {
    let equal = arguments.length > 2 && arguments[2] !== void 0 ? arguments[2] : true;
    if (a == b !== (get_proxied_value(a) == get_proxied_value(b))) {
        w.state_proxy_equality_mismatch(equal ? '==' : '!=');
    }
    return a == b === equal;
}


}),
339: (function (__unused_webpack_module, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  Tc: () => (tag),
  _e: () => (tag_proxy),
  ho: () => (tracing_expressions),
  sv: () => (get_stack)
});
/* ESM import */var _constants_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(178);
/* ESM import */var _shared_clone_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(826);
/* ESM import */var _shared_utils_js__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(986);
/* ESM import */var _client_constants__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(924);
/* ESM import */var _reactivity_effects_js__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(480);
/* ESM import */var _runtime_js__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(513);
/** @import { Derived, Reaction, Value } from '#client' */ 





/**
 * @typedef {{
 *   traces: Error[];
 * }} TraceEntry
 */ /** @type {{ reaction: Reaction | null, entries: Map<Value, TraceEntry> } | null} */ let tracing_expressions = null;
/**
 * @param {Value} signal
 * @param {TraceEntry} [entry]
 */ function log_entry(signal, entry) {
    const value = signal.v;
    if (value === UNINITIALIZED) {
        return;
    }
    const type = get_type(signal);
    const current_reaction = /** @type {Reaction} */ active_reaction;
    const dirty = signal.wv > current_reaction.wv || current_reaction.wv === 0;
    const style = dirty ? 'color: CornflowerBlue; font-weight: bold' : 'color: grey; font-weight: normal';
    // eslint-disable-next-line no-console
    console.groupCollapsed(signal.label ? `%c${type}%c ${signal.label}` : `%c${type}%c`, style, dirty ? 'font-weight: normal' : style, typeof value === 'object' && value !== null && STATE_SYMBOL in value ? snapshot(value, true) : value);
    if (type === '$derived') {
        const deps = new Set(/** @type {Derived} */ signal.deps);
        for (const dep of deps){
            log_entry(dep);
        }
    }
    if (signal.created) {
        // eslint-disable-next-line no-console
        console.log(signal.created);
    }
    if (dirty && signal.updated) {
        for (const updated of signal.updated.values()){
            // eslint-disable-next-line no-console
            console.log(updated.error);
        }
    }
    if (entry) {
        for (var trace of entry.traces){
            // eslint-disable-next-line no-console
            console.log(trace);
        }
    }
    // eslint-disable-next-line no-console
    console.groupEnd();
}
/**
 * @param {Value} signal
 * @returns {'$state' | '$derived' | 'store'}
 */ function get_type(signal) {
    var _signal_label;
    if ((signal.f & (DERIVED | ASYNC)) !== 0) return '$derived';
    return ((_signal_label = signal.label) === null || _signal_label === void 0 ? void 0 : _signal_label.startsWith('$')) ? 'store' : '$state';
}
/**
 * @template T
 * @param {() => string} label
 * @param {() => T} fn
 */ function trace(label, fn) {
    var previously_tracing_expressions = tracing_expressions;
    try {
        tracing_expressions = {
            entries: new Map(),
            reaction: active_reaction
        };
        var start = performance.now();
        var value = fn();
        var time = (performance.now() - start).toFixed(2);
        var prefix = untrack(label);
        if (!effect_tracking()) {
            // eslint-disable-next-line no-console
            console.log(`${prefix} %cran outside of an effect (${time}ms)`, 'color: grey');
        } else if (tracing_expressions.entries.size === 0) {
            // eslint-disable-next-line no-console
            console.log(`${prefix} %cno reactive dependencies (${time}ms)`, 'color: grey');
        } else {
            // eslint-disable-next-line no-console
            console.group(`${prefix} %c(${time}ms)`, 'color: grey');
            var entries = tracing_expressions.entries;
            untrack(()=>{
                for (const [signal, traces] of entries){
                    log_entry(signal, traces);
                }
            });
            tracing_expressions = null;
            // eslint-disable-next-line no-console
            console.groupEnd();
        }
        return value;
    } finally{
        tracing_expressions = previously_tracing_expressions;
    }
}
/**
 * @param {string} label
 * @returns {Error & { stack: string } | null}
 */ function get_stack(label) {
    let error = Error();
    const stack = error.stack;
    if (!stack) return null;
    const lines = stack.split('\n');
    const new_lines = [
        '\n'
    ];
    for(let i = 0; i < lines.length; i++){
        const line = lines[i];
        if (line === 'Error') {
            continue;
        }
        if (line.includes('validate_each_keys')) {
            return null;
        }
        if (line.includes('svelte/src/internal')) {
            continue;
        }
        new_lines.push(line);
    }
    if (new_lines.length === 1) {
        return null;
    }
    (0,_shared_utils_js__WEBPACK_IMPORTED_MODULE_2__/* .define_property */.Qu)(error, 'stack', {
        value: new_lines.join('\n')
    });
    (0,_shared_utils_js__WEBPACK_IMPORTED_MODULE_2__/* .define_property */.Qu)(error, 'name', {
        // 'Error' suffix is required for stack traces to be rendered properly
        value: `${label}Error`
    });
    return /** @type {Error & { stack: string }} */ error;
}
/**
 * @param {Value} source
 * @param {string} label
 */ function tag(source, label) {
    source.label = label;
    tag_proxy(source.v, label);
    return source;
}
/**
 * @param {unknown} value
 * @param {string} label
 */ function tag_proxy(value, label) {
    var // @ts-expect-error
    _value_PROXY_PATH_SYMBOL;
    value === null || value === void 0 ? void 0 : (_value_PROXY_PATH_SYMBOL = value[_client_constants__WEBPACK_IMPORTED_MODULE_3__/* .PROXY_PATH_SYMBOL */.Qf]) === null || _value_PROXY_PATH_SYMBOL === void 0 ? void 0 : _value_PROXY_PATH_SYMBOL.call(value, label);
    return value;
}
/**
 * @param {unknown} value
 */ function label(value) {
    if (typeof value === 'symbol') return `Symbol(${value.description})`;
    if (typeof value === 'function') return '<function>';
    if (typeof value === 'object' && value) return '<object>';
    return String(value);
}


}),
899: (function (__unused_webpack_module, __webpack_exports__, __webpack_require__) {

// EXPORTS
__webpack_require__.d(__webpack_exports__, {
  pP: () => (/* binding */ boundary_boundary)
});

// UNUSED EXPORTS: Boundary, pending, get_boundary

// EXTERNAL MODULE: ./node_modules/.pnpm/@swc+helpers@0.5.17/node_modules/@swc/helpers/esm/_class_private_field_get.js + 1 modules
var _class_private_field_get = __webpack_require__(570);
// EXTERNAL MODULE: ./node_modules/.pnpm/@swc+helpers@0.5.17/node_modules/@swc/helpers/esm/_class_private_field_init.js
var _class_private_field_init = __webpack_require__(636);
// EXTERNAL MODULE: ./node_modules/.pnpm/@swc+helpers@0.5.17/node_modules/@swc/helpers/esm/_class_private_field_set.js + 1 modules
var _class_private_field_set = __webpack_require__(549);
// EXTERNAL MODULE: ./node_modules/.pnpm/@swc+helpers@0.5.17/node_modules/@swc/helpers/esm/_class_private_method_get.js
var _class_private_method_get = __webpack_require__(585);
// EXTERNAL MODULE: ./node_modules/.pnpm/@swc+helpers@0.5.17/node_modules/@swc/helpers/esm/_class_private_method_init.js
var _class_private_method_init = __webpack_require__(23);
// EXTERNAL MODULE: ./node_modules/.pnpm/@swc+helpers@0.5.17/node_modules/@swc/helpers/esm/_define_property.js
var _define_property = __webpack_require__(925);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/constants.js
var constants = __webpack_require__(924);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/constants.js
var src_constants = __webpack_require__(178);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/context.js
var context = __webpack_require__(754);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/error-handling.js
var error_handling = __webpack_require__(621);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/reactivity/effects.js
var effects = __webpack_require__(480);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/runtime.js
var runtime = __webpack_require__(513);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/hydration.js
var hydration = __webpack_require__(452);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/operations.js
var operations = __webpack_require__(518);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/task.js
var task = __webpack_require__(593);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/errors.js
var errors = __webpack_require__(626);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/warnings.js
var warnings = __webpack_require__(32);
// EXTERNAL MODULE: ./node_modules/.pnpm/esm-env@1.2.2/node_modules/esm-env/false.js
var esm_env_false = __webpack_require__(832);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/reactivity/batch.js
var batch = __webpack_require__(410);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/reactivity/sources.js
var sources = __webpack_require__(264);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dev/tracing.js
var tracing = __webpack_require__(339);
;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/reactivity/create-subscriber.js






/**
 * Returns a `subscribe` function that integrates external event-based systems with Svelte's reactivity.
 * It's particularly useful for integrating with web APIs like `MediaQuery`, `IntersectionObserver`, or `WebSocket`.
 *
 * If `subscribe` is called inside an effect (including indirectly, for example inside a getter),
 * the `start` callback will be called with an `update` function. Whenever `update` is called, the effect re-runs.
 *
 * If `start` returns a cleanup function, it will be called when the effect is destroyed.
 *
 * If `subscribe` is called in multiple effects, `start` will only be called once as long as the effects
 * are active, and the returned teardown function will only be called when all effects are destroyed.
 *
 * It's best understood with an example. Here's an implementation of [`MediaQuery`](https://svelte.dev/docs/svelte/svelte-reactivity#MediaQuery):
 *
 * ```js
 * import { createSubscriber } from 'svelte/reactivity';
 * import { on } from 'svelte/events';
 *
 * export class MediaQuery {
 * 	#query;
 * 	#subscribe;
 *
 * 	constructor(query) {
 * 		this.#query = window.matchMedia(`(${query})`);
 *
 * 		this.#subscribe = createSubscriber((update) => {
 * 			// when the `change` event occurs, re-run any effects that read `this.current`
 * 			const off = on(this.#query, 'change', update);
 *
 * 			// stop listening when all the effects are destroyed
 * 			return () => off();
 * 		});
 * 	}
 *
 * 	get current() {
 * 		// This makes the getter reactive, if read in an effect
 * 		this.#subscribe();
 *
 * 		// Return the current state of the query, whether or not we're in an effect
 * 		return this.#query.matches;
 * 	}
 * }
 * ```
 * @param {(update: () => void) => (() => void) | void} start
 * @since 5.7.0
 */ function createSubscriber(start) {
    let subscribers = 0;
    let version = (0,sources/* .source */.sP)(0);
    /** @type {(() => void) | void} */ let stop;
    if (esm_env_false/* ["default"] */.A) {
        (0,tracing/* .tag */.Tc)(version, 'createSubscriber version');
    }
    return ()=>{
        if ((0,effects/* .effect_tracking */.oJ)()) {
            (0,runtime/* .get */.Jt)(version);
            (0,effects/* .render_effect */.VB)(()=>{
                if (subscribers === 0) {
                    stop = (0,runtime/* .untrack */.vz)(()=>start(()=>(0,sources/* .increment */.GV)(version)));
                }
                subscribers += 1;
                return ()=>{
                    (0,task/* .queue_micro_task */.$r)(()=>{
                        // Only count down after a microtask, else we would reach 0 before our own render effect reruns,
                        // but reach 1 again when the tick callback of the prior teardown runs. That would mean we
                        // re-subcribe unnecessarily and create a memory leak because the old subscription is never cleaned up.
                        subscribers -= 1;
                        if (subscribers === 0) {
                            stop === null || stop === void 0 ? void 0 : stop();
                            stop = undefined;
                            // Increment the version to ensure any dependent deriveds are marked dirty when the subscription is picked up again later.
                            // If we didn't do this then the comparison of write versions would determine that the derived has a later version than
                            // the subscriber, and it would not be re-run.
                            (0,sources/* .increment */.GV)(version);
                        }
                    });
                };
            });
        }
    };
}

;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/blocks/boundary.js
/** @import { Effect, Source, TemplateNode, } from '#client' */ 





















/**
 * @typedef {{
 * 	 onerror?: (error: unknown, reset: () => void) => void;
 *   failed?: (anchor: Node, error: () => unknown, reset: () => () => void) => void;
 *   pending?: (anchor: Node) => void;
 * }} BoundaryProps
 */ var flags = constants/* .EFFECT_TRANSPARENT */.lQ | constants/* .EFFECT_PRESERVED */.V$ | constants/* .BOUNDARY_EFFECT */.bp;
/**
 * @param {TemplateNode} node
 * @param {BoundaryProps} props
 * @param {((anchor: Node) => void)} children
 * @returns {void}
 */ function boundary_boundary(node, props, children) {
    new Boundary(node, props, children);
}
var _pending = /*#__PURE__*/ new WeakMap(), /** @type {TemplateNode} */ _anchor = /*#__PURE__*/ new WeakMap(), /** @type {TemplateNode | null} */ _hydrate_open = /*#__PURE__*/ new WeakMap(), /** @type {BoundaryProps} */ _props = /*#__PURE__*/ new WeakMap(), /** @type {((anchor: Node) => void)} */ _children = /*#__PURE__*/ new WeakMap(), /** @type {Effect} */ _effect = /*#__PURE__*/ new WeakMap(), /** @type {Effect | null} */ _main_effect = /*#__PURE__*/ new WeakMap(), /** @type {Effect | null} */ _pending_effect = /*#__PURE__*/ new WeakMap(), /** @type {Effect | null} */ _failed_effect = /*#__PURE__*/ new WeakMap(), /** @type {DocumentFragment | null} */ _offscreen_fragment = /*#__PURE__*/ new WeakMap(), _local_pending_count = /*#__PURE__*/ new WeakMap(), _pending_count = /*#__PURE__*/ new WeakMap(), _is_creating_fallback = /*#__PURE__*/ new WeakMap(), /**
	 * A source containing the number of pending async deriveds/expressions.
	 * Only created if `$effect.pending()` is used inside the boundary,
	 * otherwise updating the source results in needless `Batch.ensure()`
	 * calls followed by no-op flushes
	 * @type {Source<number> | null}
	 */ _effect_pending = /*#__PURE__*/ new WeakMap(), _effect_pending_update = /*#__PURE__*/ new WeakMap(), _effect_pending_subscriber = /*#__PURE__*/ new WeakMap(), _hydrate_resolved_content = /*#__PURE__*/ new WeakSet(), _hydrate_pending_content = /*#__PURE__*/ new WeakSet(), /**
	 * @param {() => Effect | null} fn
	 */ _run = /*#__PURE__*/ new WeakSet(), _show_pending_snippet = /*#__PURE__*/ new WeakSet(), /**
	 * Updates the pending count associated with the currently visible pending snippet,
	 * if any, such that we can replace the snippet with content once work is done
	 * @param {1 | -1} d
	 */ _update_pending_count = /*#__PURE__*/ new WeakSet();
class Boundary {
    /**
	 * Returns `true` if the effect exists inside a boundary whose pending snippet is shown
	 * @returns {boolean}
	 */ is_pending() {
        return (0,_class_private_field_get._)(this, _pending) || !!this.parent && this.parent.is_pending();
    }
    has_pending_snippet() {
        return !!(0,_class_private_field_get._)(this, _props).pending;
    }
    /**
	 * Update the source that powers `$effect.pending()` inside this boundary,
	 * and controls when the current `pending` snippet (if any) is removed.
	 * Do not call from inside the class
	 * @param {1 | -1} d
	 */ update_pending_count(d) {
        (0,_class_private_method_get._)(this, _update_pending_count, update_pending_count).call(this, d);
        (0,_class_private_field_set._)(this, _local_pending_count, (0,_class_private_field_get._)(this, _local_pending_count) + d);
        batch/* .effect_pending_updates.add */.x8.add((0,_class_private_field_get._)(this, _effect_pending_update));
    }
    get_effect_pending() {
        (0,_class_private_field_get._)(this, _effect_pending_subscriber).call(this);
        return (0,runtime/* .get */.Jt)((0,_class_private_field_get._)(/** @type {Source<number>} */ this, _effect_pending));
    }
    /** @param {unknown} error */ error(error) {
        var onerror = (0,_class_private_field_get._)(this, _props).onerror;
        let failed = (0,_class_private_field_get._)(this, _props).failed;
        // If we have nothing to capture the error, or if we hit an error while
        // rendering the fallback, re-throw for another boundary to handle
        if ((0,_class_private_field_get._)(this, _is_creating_fallback) || !onerror && !failed) {
            throw error;
        }
        if ((0,_class_private_field_get._)(this, _main_effect)) {
            (0,effects/* .destroy_effect */.DI)((0,_class_private_field_get._)(this, _main_effect));
            (0,_class_private_field_set._)(this, _main_effect, null);
        }
        if ((0,_class_private_field_get._)(this, _pending_effect)) {
            (0,effects/* .destroy_effect */.DI)((0,_class_private_field_get._)(this, _pending_effect));
            (0,_class_private_field_set._)(this, _pending_effect, null);
        }
        if ((0,_class_private_field_get._)(this, _failed_effect)) {
            (0,effects/* .destroy_effect */.DI)((0,_class_private_field_get._)(this, _failed_effect));
            (0,_class_private_field_set._)(this, _failed_effect, null);
        }
        if (hydration/* .hydrating */.fE) {
            (0,hydration/* .set_hydrate_node */.W0)((0,_class_private_field_get._)(/** @type {TemplateNode} */ this, _hydrate_open));
            (0,hydration/* .next */.K2)();
            (0,hydration/* .set_hydrate_node */.W0)((0,hydration/* .skip_nodes */.Ub)());
        }
        var did_reset = false;
        var calling_on_error = false;
        const reset = ()=>{
            if (did_reset) {
                warnings/* .svelte_boundary_reset_noop */.CF();
                return;
            }
            did_reset = true;
            if (calling_on_error) {
                errors/* .svelte_boundary_reset_onerror */.JJ();
            }
            // If the failure happened while flushing effects, current_batch can be null
            batch/* .Batch.ensure */.lP.ensure();
            (0,_class_private_field_set._)(this, _local_pending_count, 0);
            if ((0,_class_private_field_get._)(this, _failed_effect) !== null) {
                (0,effects/* .pause_effect */.r4)((0,_class_private_field_get._)(this, _failed_effect), ()=>{
                    (0,_class_private_field_set._)(this, _failed_effect, null);
                });
            }
            // we intentionally do not try to find the nearest pending boundary. If this boundary has one, we'll render it on reset
            // but it would be really weird to show the parent's boundary on a child reset.
            (0,_class_private_field_set._)(this, _pending, this.has_pending_snippet());
            (0,_class_private_field_set._)(this, _main_effect, (0,_class_private_method_get._)(this, _run, run).call(this, ()=>{
                (0,_class_private_field_set._)(this, _is_creating_fallback, false);
                return (0,effects/* .branch */.tk)(()=>(0,_class_private_field_get._)(this, _children).call(this, (0,_class_private_field_get._)(this, _anchor)));
            }));
            if ((0,_class_private_field_get._)(this, _pending_count) > 0) {
                (0,_class_private_method_get._)(this, _show_pending_snippet, show_pending_snippet).call(this);
            } else {
                (0,_class_private_field_set._)(this, _pending, false);
            }
        };
        var previous_reaction = runtime/* .active_reaction */.hp;
        try {
            (0,runtime/* .set_active_reaction */.G0)(null);
            calling_on_error = true;
            onerror === null || onerror === void 0 ? void 0 : onerror(error, reset);
            calling_on_error = false;
        } catch (error) {
            (0,error_handling/* .invoke_error_boundary */.n)(error, (0,_class_private_field_get._)(this, _effect) && (0,_class_private_field_get._)(this, _effect).parent);
        } finally{
            (0,runtime/* .set_active_reaction */.G0)(previous_reaction);
        }
        if (failed) {
            (0,task/* .queue_micro_task */.$r)(()=>{
                (0,_class_private_field_set._)(this, _failed_effect, (0,_class_private_method_get._)(this, _run, run).call(this, ()=>{
                    (0,_class_private_field_set._)(this, _is_creating_fallback, true);
                    try {
                        return (0,effects/* .branch */.tk)(()=>{
                            failed((0,_class_private_field_get._)(this, _anchor), ()=>error, ()=>reset);
                        });
                    } catch (error) {
                        (0,error_handling/* .invoke_error_boundary */.n)(error, /** @type {Effect} */ (0,_class_private_field_get._)(this, _effect).parent);
                        return null;
                    } finally{
                        (0,_class_private_field_set._)(this, _is_creating_fallback, false);
                    }
                }));
            });
        }
    }
    /**
	 * @param {TemplateNode} node
	 * @param {BoundaryProps} props
	 * @param {((anchor: Node) => void)} children
	 */ constructor(node, props, children){
        (0,_class_private_method_init._)(this, _hydrate_resolved_content);
        (0,_class_private_method_init._)(this, _hydrate_pending_content);
        (0,_class_private_method_init._)(this, _run);
        (0,_class_private_method_init._)(this, _show_pending_snippet);
        (0,_class_private_method_init._)(this, _update_pending_count);
        /** @type {Boundary | null} */ (0,_define_property._)(this, "parent", void 0);
        (0,_class_private_field_init._)(this, _pending, {
            writable: true,
            value: false
        });
        (0,_class_private_field_init._)(this, _anchor, {
            writable: true,
            value: void 0
        });
        (0,_class_private_field_init._)(this, _hydrate_open, {
            writable: true,
            value: hydration/* .hydrating */.fE ? hydration/* .hydrate_node */.Xb : null
        });
        (0,_class_private_field_init._)(this, _props, {
            writable: true,
            value: void 0
        });
        (0,_class_private_field_init._)(this, _children, {
            writable: true,
            value: void 0
        });
        (0,_class_private_field_init._)(this, _effect, {
            writable: true,
            value: void 0
        });
        (0,_class_private_field_init._)(this, _main_effect, {
            writable: true,
            value: null
        });
        (0,_class_private_field_init._)(this, _pending_effect, {
            writable: true,
            value: null
        });
        (0,_class_private_field_init._)(this, _failed_effect, {
            writable: true,
            value: null
        });
        (0,_class_private_field_init._)(this, _offscreen_fragment, {
            writable: true,
            value: null
        });
        (0,_class_private_field_init._)(this, _local_pending_count, {
            writable: true,
            value: 0
        });
        (0,_class_private_field_init._)(this, _pending_count, {
            writable: true,
            value: 0
        });
        (0,_class_private_field_init._)(this, _is_creating_fallback, {
            writable: true,
            value: false
        });
        (0,_class_private_field_init._)(this, _effect_pending, {
            writable: true,
            value: null
        });
        (0,_class_private_field_init._)(this, _effect_pending_update, {
            writable: true,
            value: ()=>{
                if ((0,_class_private_field_get._)(this, _effect_pending)) {
                    (0,sources/* .internal_set */.LY)((0,_class_private_field_get._)(this, _effect_pending), (0,_class_private_field_get._)(this, _local_pending_count));
                }
            }
        });
        (0,_class_private_field_init._)(this, _effect_pending_subscriber, {
            writable: true,
            value: createSubscriber(()=>{
                (0,_class_private_field_set._)(this, _effect_pending, (0,sources/* .source */.sP)((0,_class_private_field_get._)(this, _local_pending_count)));
                if (esm_env_false/* ["default"] */.A) {
                    (0,tracing/* .tag */.Tc)((0,_class_private_field_get._)(this, _effect_pending), '$effect.pending()');
                }
                return ()=>{
                    (0,_class_private_field_set._)(this, _effect_pending, null);
                };
            })
        });
        (0,_class_private_field_set._)(this, _anchor, node);
        (0,_class_private_field_set._)(this, _props, props);
        (0,_class_private_field_set._)(this, _children, children);
        this.parent = /** @type {Effect} */ runtime/* .active_effect.b */.Fg.b;
        (0,_class_private_field_set._)(this, _pending, !!(0,_class_private_field_get._)(this, _props).pending);
        (0,_class_private_field_set._)(this, _effect, (0,effects/* .block */.om)(()=>{
            /** @type {Effect} */ runtime/* .active_effect.b */.Fg.b = this;
            if (hydration/* .hydrating */.fE) {
                const comment = (0,_class_private_field_get._)(this, _hydrate_open);
                (0,hydration/* .hydrate_next */.E$)();
                const server_rendered_pending = /** @type {Comment} */ comment.nodeType === constants/* .COMMENT_NODE */.dz && /** @type {Comment} */ comment.data === src_constants/* .HYDRATION_START_ELSE */.qn;
                if (server_rendered_pending) {
                    (0,_class_private_method_get._)(this, _hydrate_pending_content, hydrate_pending_content).call(this);
                } else {
                    (0,_class_private_method_get._)(this, _hydrate_resolved_content, hydrate_resolved_content).call(this);
                }
            } else {
                try {
                    (0,_class_private_field_set._)(this, _main_effect, (0,effects/* .branch */.tk)(()=>children((0,_class_private_field_get._)(this, _anchor))));
                } catch (error) {
                    this.error(error);
                }
                if ((0,_class_private_field_get._)(this, _pending_count) > 0) {
                    (0,_class_private_method_get._)(this, _show_pending_snippet, show_pending_snippet).call(this);
                } else {
                    (0,_class_private_field_set._)(this, _pending, false);
                }
            }
        }, flags));
        if (hydration/* .hydrating */.fE) {
            (0,_class_private_field_set._)(this, _anchor, hydration/* .hydrate_node */.Xb);
        }
    }
}
function hydrate_resolved_content() {
    try {
        (0,_class_private_field_set._)(this, _main_effect, (0,effects/* .branch */.tk)(()=>(0,_class_private_field_get._)(this, _children).call(this, (0,_class_private_field_get._)(this, _anchor))));
    } catch (error) {
        this.error(error);
    }
    // Since server rendered resolved content, we never show pending state
    // Even if client-side async operations are still running, the content is already displayed
    (0,_class_private_field_set._)(this, _pending, false);
}
function hydrate_pending_content() {
    const pending = (0,_class_private_field_get._)(this, _props).pending;
    if (!pending) {
        return;
    }
    (0,_class_private_field_set._)(this, _pending_effect, (0,effects/* .branch */.tk)(()=>pending((0,_class_private_field_get._)(this, _anchor))));
    batch/* .Batch.enqueue */.lP.enqueue(()=>{
        (0,_class_private_field_set._)(this, _main_effect, (0,_class_private_method_get._)(this, _run, run).call(this, ()=>{
            batch/* .Batch.ensure */.lP.ensure();
            return (0,effects/* .branch */.tk)(()=>(0,_class_private_field_get._)(this, _children).call(this, (0,_class_private_field_get._)(this, _anchor)));
        }));
        if ((0,_class_private_field_get._)(this, _pending_count) > 0) {
            (0,_class_private_method_get._)(this, _show_pending_snippet, show_pending_snippet).call(this);
        } else {
            (0,effects/* .pause_effect */.r4)((0,_class_private_field_get._)(/** @type {Effect} */ this, _pending_effect), ()=>{
                (0,_class_private_field_set._)(this, _pending_effect, null);
            });
            (0,_class_private_field_set._)(this, _pending, false);
        }
    });
}
function run(fn) {
    var previous_effect = runtime/* .active_effect */.Fg;
    var previous_reaction = runtime/* .active_reaction */.hp;
    var previous_ctx = context/* .component_context */.UL;
    (0,runtime/* .set_active_effect */.gU)((0,_class_private_field_get._)(this, _effect));
    (0,runtime/* .set_active_reaction */.G0)((0,_class_private_field_get._)(this, _effect));
    (0,context/* .set_component_context */.De)((0,_class_private_field_get._)(this, _effect).ctx);
    try {
        return fn();
    } catch (e) {
        (0,error_handling/* .handle_error */.i)(e);
        return null;
    } finally{
        (0,runtime/* .set_active_effect */.gU)(previous_effect);
        (0,runtime/* .set_active_reaction */.G0)(previous_reaction);
        (0,context/* .set_component_context */.De)(previous_ctx);
    }
}
function show_pending_snippet() {
    const pending = /** @type {(anchor: Node) => void} */ (0,_class_private_field_get._)(this, _props).pending;
    if ((0,_class_private_field_get._)(this, _main_effect) !== null) {
        (0,_class_private_field_set._)(this, _offscreen_fragment, document.createDocumentFragment());
        move_effect((0,_class_private_field_get._)(this, _main_effect), (0,_class_private_field_get._)(this, _offscreen_fragment));
    }
    if ((0,_class_private_field_get._)(this, _pending_effect) === null) {
        (0,_class_private_field_set._)(this, _pending_effect, (0,effects/* .branch */.tk)(()=>pending((0,_class_private_field_get._)(this, _anchor))));
    }
}
function update_pending_count(d) {
    var _this_parent;
    if (!this.has_pending_snippet()) {
        if (this.parent) {
            (0,_class_private_method_get._)(_this_parent = this.parent, _update_pending_count, update_pending_count).call(_this_parent, d);
        }
        // if there's no parent, we're in a scope with no pending snippet
        return;
    }
    (0,_class_private_field_set._)(this, _pending_count, (0,_class_private_field_get._)(this, _pending_count) + d);
    if ((0,_class_private_field_get._)(this, _pending_count) === 0) {
        (0,_class_private_field_set._)(this, _pending, false);
        if ((0,_class_private_field_get._)(this, _pending_effect)) {
            (0,effects/* .pause_effect */.r4)((0,_class_private_field_get._)(this, _pending_effect), ()=>{
                (0,_class_private_field_set._)(this, _pending_effect, null);
            });
        }
        if ((0,_class_private_field_get._)(this, _offscreen_fragment)) {
            (0,_class_private_field_get._)(this, _anchor).before((0,_class_private_field_get._)(this, _offscreen_fragment));
            (0,_class_private_field_set._)(this, _offscreen_fragment, null);
        }
    }
}
/**
 *
 * @param {Effect} effect
 * @param {DocumentFragment} fragment
 */ function move_effect(effect, fragment) {
    var node = effect.nodes_start;
    var end = effect.nodes_end;
    while(node !== null){
        /** @type {TemplateNode | null} */ var next = node === end ? null : /** @type {TemplateNode} */ (0,operations/* .get_next_sibling */.M$)(node);
        fragment.append(node);
        node = next;
    }
}
function get_boundary() {
    return /** @type {Boundary} */ /** @type {Effect} */ active_effect.b;
}
function boundary_pending() {
    if (active_effect === null) {
        e.effect_pending_outside_reaction();
    }
    var boundary = active_effect.b;
    if (boundary === null) {
        return 0; // TODO eventually we will need this to be global
    }
    return boundary.get_effect_pending();
}


}),
768: (function (__unused_webpack_module, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  UA: () => (snippet)
});
/* ESM import */var _client_constants__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(924);
/* ESM import */var _reactivity_effects_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(480);
/* ESM import */var _context_js__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(754);
/* ESM import */var _hydration_js__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(452);
/* ESM import */var _template_js__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(782);
/* ESM import */var _errors_js__WEBPACK_IMPORTED_MODULE_9__ = __webpack_require__(626);
/* ESM import */var esm_env__WEBPACK_IMPORTED_MODULE_8__ = __webpack_require__(832);
/* ESM import */var _operations_js__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(518);
/* ESM import */var _shared_utils_js__WEBPACK_IMPORTED_MODULE_6__ = __webpack_require__(986);
/* ESM import */var _shared_validate_js__WEBPACK_IMPORTED_MODULE_7__ = __webpack_require__(461);
/** @import { Snippet } from 'svelte' */ /** @import { Effect, TemplateNode } from '#client' */ /** @import { Getters } from '#shared' */ 











/**
 * @template {(node: TemplateNode, ...args: any[]) => void} SnippetFn
 * @param {TemplateNode} node
 * @param {() => SnippetFn | null | undefined} get_snippet
 * @param {(() => any)[]} args
 * @returns {void}
 */ function snippet(node, get_snippet) {
    for(var _len = arguments.length, args = new Array(_len > 2 ? _len - 2 : 0), _key = 2; _key < _len; _key++){
        args[_key - 2] = arguments[_key];
    }
    var anchor = node;
    /** @type {SnippetFn | null | undefined} */ // @ts-ignore
    var snippet = _shared_utils_js__WEBPACK_IMPORTED_MODULE_6__/* .noop */.lQ;
    /** @type {Effect | null} */ var snippet_effect;
    (0,_reactivity_effects_js__WEBPACK_IMPORTED_MODULE_1__/* .block */.om)(()=>{
        if (snippet === (snippet = get_snippet())) return;
        if (snippet_effect) {
            (0,_reactivity_effects_js__WEBPACK_IMPORTED_MODULE_1__/* .destroy_effect */.DI)(snippet_effect);
            snippet_effect = null;
        }
        if (esm_env__WEBPACK_IMPORTED_MODULE_8__/* ["default"] */.A && snippet == null) {
            _errors_js__WEBPACK_IMPORTED_MODULE_9__/* .invalid_snippet */.WR();
        }
        snippet_effect = (0,_reactivity_effects_js__WEBPACK_IMPORTED_MODULE_1__/* .branch */.tk)(()=>/** @type {SnippetFn} */ snippet(anchor, ...args));
    }, _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .EFFECT_TRANSPARENT */.lQ);
    if (_hydration_js__WEBPACK_IMPORTED_MODULE_3__/* .hydrating */.fE) {
        anchor = _hydration_js__WEBPACK_IMPORTED_MODULE_3__/* .hydrate_node */.Xb;
    }
}
/**
 * In development, wrap the snippet function so that it passes validation, and so that the
 * correct component context is set for ownership checks
 * @param {any} component
 * @param {(node: TemplateNode, ...args: any[]) => void} fn
 */ function wrap_snippet(component, fn) {
    const snippet = function(/** @type {TemplateNode} */ node) {
        for(var _len = arguments.length, args = new Array(_len > 1 ? _len - 1 : 0), _key = 1; _key < _len; _key++){
            args[_key - 1] = arguments[_key];
        }
        var previous_component_function = dev_current_component_function;
        set_dev_current_component_function(component);
        try {
            return fn(node, ...args);
        } finally{
            set_dev_current_component_function(previous_component_function);
        }
    };
    prevent_snippet_stringification(snippet);
    return snippet;
}
/**
 * Create a snippet programmatically
 * @template {unknown[]} Params
 * @param {(...params: Getters<Params>) => {
 *   render: () => string
 *   setup?: (element: Element) => void | (() => void)
 * }} fn
 * @returns {Snippet<Params>}
 */ function createRawSnippet(fn) {
    // @ts-expect-error the types are a lie
    return function(/** @type {TemplateNode} */ anchor) {
        for(var _len = arguments.length, params = new Array(_len > 1 ? _len - 1 : 0), _key = 1; _key < _len; _key++){
            params[_key - 1] = arguments[_key];
        }
        var _snippet_setup;
        var snippet = fn(...params);
        /** @type {Element} */ var element;
        if (hydrating) {
            element = /** @type {Element} */ hydrate_node;
            hydrate_next();
        } else {
            var html = snippet.render().trim();
            var fragment = create_fragment_from_html(html);
            element = /** @type {Element} */ get_first_child(fragment);
            if (DEV && (get_next_sibling(element) !== null || element.nodeType !== ELEMENT_NODE)) {
                w.invalid_raw_snippet_render();
            }
            anchor.before(element);
        }
        const result = (_snippet_setup = snippet.setup) === null || _snippet_setup === void 0 ? void 0 : _snippet_setup.call(snippet, element);
        assign_nodes(element, element);
        if (typeof result === 'function') {
            teardown(result);
        }
    };
}


}),
777: (function (__unused_webpack_module, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  j: () => (reset_head_anchor)
});
/* ESM import */var _hydration_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(452);
/* ESM import */var _operations_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(518);
/* ESM import */var _reactivity_effects_js__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(480);
/* ESM import */var _client_constants__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(924);
/* ESM import */var _constants_js__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(178);
/** @import { TemplateNode } from '#client' */ 




/**
 * @type {Node | undefined}
 */ let head_anchor;
function reset_head_anchor() {
    head_anchor = undefined;
}
/**
 * @param {(anchor: Node) => void} render_fn
 * @returns {void}
 */ function head(render_fn) {
    // The head function may be called after the first hydration pass and ssr comment nodes may still be present,
    // therefore we need to skip that when we detect that we're not in hydration mode.
    let previous_hydrate_node = null;
    let was_hydrating = hydrating;
    /** @type {Comment | Text} */ var anchor;
    if (hydrating) {
        previous_hydrate_node = hydrate_node;
        // There might be multiple head blocks in our app, so we need to account for each one needing independent hydration.
        if (head_anchor === undefined) {
            head_anchor = /** @type {TemplateNode} */ get_first_child(document.head);
        }
        while(head_anchor !== null && (head_anchor.nodeType !== COMMENT_NODE || /** @type {Comment} */ head_anchor.data !== HYDRATION_START)){
            head_anchor = /** @type {TemplateNode} */ get_next_sibling(head_anchor);
        }
        // If we can't find an opening hydration marker, skip hydration (this can happen
        // if a framework rendered body but not head content)
        if (head_anchor === null) {
            set_hydrating(false);
        } else {
            head_anchor = set_hydrate_node(/** @type {TemplateNode} */ get_next_sibling(head_anchor));
        }
    }
    if (!hydrating) {
        anchor = document.head.appendChild(create_text());
    }
    try {
        block(()=>render_fn(anchor), HEAD_EFFECT);
    } finally{
        if (was_hydrating) {
            set_hydrating(true);
            head_anchor = hydrate_node; // so that next head block starts from the correct node
            set_hydrate_node(/** @type {TemplateNode} */ previous_hydrate_node);
        }
    }
}


}),
408: (function (__unused_webpack_module, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  $w: () => (without_reactive_context)
});
/* ESM import */var _reactivity_effects_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(480);
/* ESM import */var _runtime_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(513);
/* ESM import */var _misc_js__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(108);



/**
 * Fires the handler once immediately (unless corresponding arg is set to `false`),
 * then listens to the given events until the render effect context is destroyed
 * @param {EventTarget} target
 * @param {Array<string>} events
 * @param {(event?: Event) => void} handler
 * @param {any} call_handler_immediately
 */ function listen(target, events, handler) {
    let call_handler_immediately = arguments.length > 3 && arguments[3] !== void 0 ? arguments[3] : true;
    if (call_handler_immediately) {
        handler();
    }
    for (var name of events){
        target.addEventListener(name, handler);
    }
    teardown(()=>{
        for (var name of events){
            target.removeEventListener(name, handler);
        }
    });
}
/**
 * @template T
 * @param {() => T} fn
 */ function without_reactive_context(fn) {
    var previous_reaction = _runtime_js__WEBPACK_IMPORTED_MODULE_1__/* .active_reaction */.hp;
    var previous_effect = _runtime_js__WEBPACK_IMPORTED_MODULE_1__/* .active_effect */.Fg;
    (0,_runtime_js__WEBPACK_IMPORTED_MODULE_1__/* .set_active_reaction */.G0)(null);
    (0,_runtime_js__WEBPACK_IMPORTED_MODULE_1__/* .set_active_effect */.gU)(null);
    try {
        return fn();
    } finally{
        (0,_runtime_js__WEBPACK_IMPORTED_MODULE_1__/* .set_active_reaction */.G0)(previous_reaction);
        (0,_runtime_js__WEBPACK_IMPORTED_MODULE_1__/* .set_active_effect */.gU)(previous_effect);
    }
}
/**
 * Listen to the given event, and then instantiate a global form reset listener if not already done,
 * to notify all bindings when the form is reset
 * @param {HTMLElement} element
 * @param {string} event
 * @param {(is_reset?: true) => void} handler
 * @param {(is_reset?: true) => void} [on_reset]
 */ function listen_to_event_and_reset_event(element, event, handler) {
    let on_reset = arguments.length > 3 && arguments[3] !== void 0 ? arguments[3] : handler;
    element.addEventListener(event, ()=>without_reactive_context(handler));
    // @ts-expect-error
    const prev = element.__on_r;
    if (prev) {
        // special case for checkbox that can have multiple binds (group & checked)
        // @ts-expect-error
        element.__on_r = ()=>{
            prev();
            on_reset(true);
        };
    } else {
        // @ts-expect-error
        element.__on_r = ()=>on_reset(true);
    }
    add_form_reset_listener();
}


}),
417: (function (__unused_webpack_module, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  Mm: () => (delegate),
  Sr: () => (root_event_handles),
  Ts: () => (all_registered_events),
  n7: () => (handle_event_propagation)
});
/* ESM import */var _reactivity_effects_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(480);
/* ESM import */var _shared_utils_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(986);
/* ESM import */var _hydration_js__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(452);
/* ESM import */var _task_js__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(593);
/* ESM import */var _constants_js__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(178);
/* ESM import */var _runtime_js__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(513);
/* ESM import */var _bindings_shared_js__WEBPACK_IMPORTED_MODULE_6__ = __webpack_require__(408);








/** @type {Set<string>} */ const all_registered_events = new Set();
/** @type {Set<(events: Array<string>) => void>} */ const root_event_handles = new Set();
/**
 * SSR adds onload and onerror attributes to catch those events before the hydration.
 * This function detects those cases, removes the attributes and replays the events.
 * @param {HTMLElement} dom
 */ function replay_events(dom) {
    if (!hydrating) return;
    dom.removeAttribute('onload');
    dom.removeAttribute('onerror');
    // @ts-expect-error
    const event = dom.__e;
    if (event !== undefined) {
        // @ts-expect-error
        dom.__e = undefined;
        queueMicrotask(()=>{
            if (dom.isConnected) {
                dom.dispatchEvent(event);
            }
        });
    }
}
/**
 * @param {string} event_name
 * @param {EventTarget} dom
 * @param {EventListener} [handler]
 * @param {AddEventListenerOptions} [options]
 */ function create_event(event_name, dom, handler) {
    let options = arguments.length > 3 && arguments[3] !== void 0 ? arguments[3] : {};
    /**
	 * @this {EventTarget}
	 */ function target_handler(/** @type {Event} */ event) {
        if (!options.capture) {
            // Only call in the bubble phase, else delegated events would be called before the capturing events
            handle_event_propagation.call(dom, event);
        }
        if (!event.cancelBubble) {
            return without_reactive_context(()=>{
                return handler === null || handler === void 0 ? void 0 : handler.call(this, event);
            });
        }
    }
    // Chrome has a bug where pointer events don't work when attached to a DOM element that has been cloned
    // with cloneNode() and the DOM element is disconnected from the document. To ensure the event works, we
    // defer the attachment till after it's been appended to the document. TODO: remove this once Chrome fixes
    // this bug. The same applies to wheel events and touch events.
    if (event_name.startsWith('pointer') || event_name.startsWith('touch') || event_name === 'wheel') {
        queue_micro_task(()=>{
            dom.addEventListener(event_name, target_handler, options);
        });
    } else {
        dom.addEventListener(event_name, target_handler, options);
    }
    return target_handler;
}
/**
 * Attaches an event handler to an element and returns a function that removes the handler. Using this
 * rather than `addEventListener` will preserve the correct order relative to handlers added declaratively
 * (with attributes like `onclick`), which use event delegation for performance reasons
 *
 * @param {EventTarget} element
 * @param {string} type
 * @param {EventListener} handler
 * @param {AddEventListenerOptions} [options]
 */ function on(element, type, handler) {
    let options = arguments.length > 3 && arguments[3] !== void 0 ? arguments[3] : {};
    var target_handler = create_event(type, element, handler, options);
    return ()=>{
        element.removeEventListener(type, target_handler, options);
    };
}
/**
 * @param {string} event_name
 * @param {Element} dom
 * @param {EventListener} [handler]
 * @param {boolean} [capture]
 * @param {boolean} [passive]
 * @returns {void}
 */ function event(event_name, dom, handler, capture, passive) {
    var options = {
        capture,
        passive
    };
    var target_handler = create_event(event_name, dom, handler, options);
    if (dom === document.body || // @ts-ignore
    dom === window || // @ts-ignore
    dom === document || // Firefox has quirky behavior, it can happen that we still get "canplay" events when the element is already removed
    dom instanceof HTMLMediaElement) {
        teardown(()=>{
            dom.removeEventListener(event_name, target_handler, options);
        });
    }
}
/**
 * @param {Array<string>} events
 * @returns {void}
 */ function delegate(events) {
    for(var i = 0; i < events.length; i++){
        all_registered_events.add(events[i]);
    }
    for (var fn of root_event_handles){
        fn(events);
    }
}
// used to store the reference to the currently propagated event
// to prevent garbage collection between microtasks in Firefox
// If the event object is GCed too early, the expando __root property
// set on the event object is lost, causing the event delegation
// to process the event twice
let last_propagated_event = null;
/**
 * @this {EventTarget}
 * @param {Event} event
 * @returns {void}
 */ function handle_event_propagation(event) {
    var _event_composedPath;
    var handler_element = this;
    var owner_document = /** @type {Node} */ handler_element.ownerDocument;
    var event_name = event.type;
    var path = ((_event_composedPath = event.composedPath) === null || _event_composedPath === void 0 ? void 0 : _event_composedPath.call(event)) || [];
    var current_target = /** @type {null | Element} */ path[0] || event.target;
    last_propagated_event = event;
    // composedPath contains list of nodes the event has propagated through.
    // We check __root to skip all nodes below it in case this is a
    // parent of the __root node, which indicates that there's nested
    // mounted apps. In this case we don't want to trigger events multiple times.
    var path_idx = 0;
    // the `last_propagated_event === event` check is redundant, but
    // without it the variable will be DCE'd and things will
    // fail mysteriously in Firefox
    // @ts-expect-error is added below
    var handled_at = last_propagated_event === event && event.__root;
    if (handled_at) {
        var at_idx = path.indexOf(handled_at);
        if (at_idx !== -1 && (handler_element === document || handler_element === /** @type {any} */ window)) {
            // This is the fallback document listener or a window listener, but the event was already handled
            // -> ignore, but set handle_at to document/window so that we're resetting the event
            // chain in case someone manually dispatches the same event object again.
            // @ts-expect-error
            event.__root = handler_element;
            return;
        }
        // We're deliberately not skipping if the index is higher, because
        // someone could create an event programmatically and emit it multiple times,
        // in which case we want to handle the whole propagation chain properly each time.
        // (this will only be a false negative if the event is dispatched multiple times and
        // the fallback document listener isn't reached in between, but that's super rare)
        var handler_idx = path.indexOf(handler_element);
        if (handler_idx === -1) {
            // handle_idx can theoretically be -1 (happened in some JSDOM testing scenarios with an event listener on the window object)
            // so guard against that, too, and assume that everything was handled at this point.
            return;
        }
        if (at_idx <= handler_idx) {
            path_idx = at_idx;
        }
    }
    current_target = /** @type {Element} */ path[path_idx] || event.target;
    // there can only be one delegated event per element, and we either already handled the current target,
    // or this is the very first target in the chain which has a non-delegated listener, in which case it's safe
    // to handle a possible delegated event on it later (through the root delegation listener for example).
    if (current_target === handler_element) return;
    // Proxy currentTarget to correct target
    (0,_shared_utils_js__WEBPACK_IMPORTED_MODULE_1__/* .define_property */.Qu)(event, 'currentTarget', {
        configurable: true,
        get () {
            return current_target || owner_document;
        }
    });
    // This started because of Chromium issue https://chromestatus.com/feature/5128696823545856,
    // where removal or moving of of the DOM can cause sync `blur` events to fire, which can cause logic
    // to run inside the current `active_reaction`, which isn't what we want at all. However, on reflection,
    // it's probably best that all event handled by Svelte have this behaviour, as we don't really want
    // an event handler to run in the context of another reaction or effect.
    var previous_reaction = _runtime_js__WEBPACK_IMPORTED_MODULE_5__/* .active_reaction */.hp;
    var previous_effect = _runtime_js__WEBPACK_IMPORTED_MODULE_5__/* .active_effect */.Fg;
    (0,_runtime_js__WEBPACK_IMPORTED_MODULE_5__/* .set_active_reaction */.G0)(null);
    (0,_runtime_js__WEBPACK_IMPORTED_MODULE_5__/* .set_active_effect */.gU)(null);
    try {
        /**
		 * @type {unknown}
		 */ var throw_error;
        /**
		 * @type {unknown[]}
		 */ var other_errors = [];
        while(current_target !== null){
            /** @type {null | Element} */ var parent_element = current_target.assignedSlot || current_target.parentNode || /** @type {any} */ current_target.host || null;
            try {
                // @ts-expect-error
                var delegated = current_target['__' + event_name];
                if (delegated != null && (!/** @type {any} */ current_target.disabled || // DOM could've been updated already by the time this is reached, so we check this as well
                // -> the target could not have been disabled because it emits the event in the first place
                event.target === current_target)) {
                    if ((0,_shared_utils_js__WEBPACK_IMPORTED_MODULE_1__/* .is_array */.PI)(delegated)) {
                        var [fn, ...data] = delegated;
                        fn.apply(current_target, [
                            event,
                            ...data
                        ]);
                    } else {
                        delegated.call(current_target, event);
                    }
                }
            } catch (error) {
                if (throw_error) {
                    other_errors.push(error);
                } else {
                    throw_error = error;
                }
            }
            if (event.cancelBubble || parent_element === handler_element || parent_element === null) {
                break;
            }
            current_target = parent_element;
        }
        if (throw_error) {
            for (let error of other_errors){
                // Throw the rest of the errors, one-by-one on a microtask
                queueMicrotask(()=>{
                    throw error;
                });
            }
            throw throw_error;
        }
    } finally{
        // @ts-expect-error is used above
        event.__root = handler_element;
        // @ts-ignore remove proxy on currentTarget
        delete event.currentTarget;
        (0,_runtime_js__WEBPACK_IMPORTED_MODULE_5__/* .set_active_reaction */.G0)(previous_reaction);
        (0,_runtime_js__WEBPACK_IMPORTED_MODULE_5__/* .set_active_effect */.gU)(previous_effect);
    }
}
/**
 * In dev, warn if an event handler is not a function, as it means the
 * user probably called the handler or forgot to add a `() =>`
 * @param {() => (event: Event, ...args: any) => void} thunk
 * @param {EventTarget} element
 * @param {[Event, ...any]} args
 * @param {any} component
 * @param {[number, number]} [loc]
 * @param {boolean} [remove_parens]
 */ function apply(thunk, element, args, component, loc) {
    let has_side_effects = arguments.length > 5 && arguments[5] !== void 0 ? arguments[5] : false, remove_parens = arguments.length > 6 && arguments[6] !== void 0 ? arguments[6] : false;
    let handler;
    let error;
    try {
        handler = thunk();
    } catch (e) {
        error = e;
    }
    if (typeof handler !== 'function' && (has_side_effects || handler != null || error)) {
        var _args_, _args_1;
        const filename = component === null || component === void 0 ? void 0 : component[FILENAME];
        const location = loc ? ` at ${filename}:${loc[0]}:${loc[1]}` : ` in ${filename}`;
        const phase = ((_args_ = args[0]) === null || _args_ === void 0 ? void 0 : _args_.eventPhase) < Event.BUBBLING_PHASE ? 'capture' : '';
        const event_name = ((_args_1 = args[0]) === null || _args_1 === void 0 ? void 0 : _args_1.type) + phase;
        const description = `\`${event_name}\` handler${location}`;
        const suggestion = remove_parens ? 'remove the trailing `()`' : 'add a leading `() =>`';
        w.event_handler_invalid(description, suggestion);
        if (error) {
            throw error;
        }
    }
    handler === null || handler === void 0 ? void 0 : handler.apply(element, args);
}


}),
108: (function (__unused_webpack_module, __unused_webpack___webpack_exports__, __webpack_require__) {
/* ESM import */var _hydration_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(452);
/* ESM import */var _operations_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(518);
/* ESM import */var _task_js__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(593);



/**
 * @param {HTMLElement} dom
 * @param {boolean} value
 * @returns {void}
 */ function autofocus(dom, value) {
    if (value) {
        const body = document.body;
        dom.autofocus = true;
        queue_micro_task(()=>{
            if (document.activeElement === body) {
                dom.focus();
            }
        });
    }
}
/**
 * The child of a textarea actually corresponds to the defaultValue property, so we need
 * to remove it upon hydration to avoid a bug when someone resets the form value.
 * @param {HTMLTextAreaElement} dom
 * @returns {void}
 */ function remove_textarea_child(dom) {
    if (hydrating && get_first_child(dom) !== null) {
        clear_text_content(dom);
    }
}
let listening_to_form_reset = false;
function add_form_reset_listener() {
    if (!listening_to_form_reset) {
        listening_to_form_reset = true;
        document.addEventListener('reset', (evt)=>{
            // Needs to happen one tick later or else the dom properties of the form
            // elements have not updated to their reset values yet
            Promise.resolve().then(()=>{
                if (!evt.defaultPrevented) {
                    for (const e of evt.target.elements){
                        var // @ts-expect-error
                        _e___on_r;
                        (_e___on_r = e.__on_r) === null || _e___on_r === void 0 ? void 0 : _e___on_r.call(e);
                    }
                }
            });
        }, // In the capture phase to guarantee we get noticed of it (no possiblity of stopPropagation)
        {
            capture: true
        });
    }
}


}),
452: (function (__unused_webpack_module, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  E$: () => (hydrate_next),
  K2: () => (next),
  Ub: () => (skip_nodes),
  W0: () => (set_hydrate_node),
  Xb: () => (hydrate_node),
  cL: () => (reset),
  fE: () => (hydrating),
  mK: () => (set_hydrating),
  no: () => (read_hydration_instruction)
});
/* ESM import */var _client_constants__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(924);
/* ESM import */var _constants_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(178);
/* ESM import */var _warnings_js__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(32);
/* ESM import */var _operations_js__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(518);
/** @import { TemplateNode } from '#client' */ 



/**
 * Use this variable to guard everything related to hydration code so it can be treeshaken out
 * if the user doesn't use the `hydrate` method and these code paths are therefore not needed.
 */ let hydrating = false;
/** @param {boolean} value */ function set_hydrating(value) {
    hydrating = value;
}
/**
 * The node that is currently being hydrated. This starts out as the first node inside the opening
 * <!--[--> comment, and updates each time a component calls `$.child(...)` or `$.sibling(...)`.
 * When entering a block (e.g. `{#if ...}`), `hydrate_node` is the block opening comment; by the
 * time we leave the block it is the closing comment, which serves as the block's anchor.
 * @type {TemplateNode}
 */ let hydrate_node;
/** @param {TemplateNode} node */ function set_hydrate_node(node) {
    if (node === null) {
        _warnings_js__WEBPACK_IMPORTED_MODULE_3__/* .hydration_mismatch */.eZ();
        throw _constants_js__WEBPACK_IMPORTED_MODULE_1__/* .HYDRATION_ERROR */.kD;
    }
    return hydrate_node = node;
}
function hydrate_next() {
    return set_hydrate_node(/** @type {TemplateNode} */ (0,_operations_js__WEBPACK_IMPORTED_MODULE_2__/* .get_next_sibling */.M$)(hydrate_node));
}
/** @param {TemplateNode} node */ function reset(node) {
    if (!hydrating) return;
    // If the node has remaining siblings, something has gone wrong
    if ((0,_operations_js__WEBPACK_IMPORTED_MODULE_2__/* .get_next_sibling */.M$)(hydrate_node) !== null) {
        _warnings_js__WEBPACK_IMPORTED_MODULE_3__/* .hydration_mismatch */.eZ();
        throw _constants_js__WEBPACK_IMPORTED_MODULE_1__/* .HYDRATION_ERROR */.kD;
    }
    hydrate_node = node;
}
/**
 * @param {HTMLTemplateElement} template
 */ function hydrate_template(template) {
    if (hydrating) {
        // @ts-expect-error TemplateNode doesn't include DocumentFragment, but it's actually fine
        hydrate_node = template.content;
    }
}
function next() {
    let count = arguments.length > 0 && arguments[0] !== void 0 ? arguments[0] : 1;
    if (hydrating) {
        var i = count;
        var node = hydrate_node;
        while(i--){
            node = /** @type {TemplateNode} */ (0,_operations_js__WEBPACK_IMPORTED_MODULE_2__/* .get_next_sibling */.M$)(node);
        }
        hydrate_node = node;
    }
}
/**
 * Skips or removes (depending on {@link remove}) all nodes starting at `hydrate_node` up until the next hydration end comment
 * @param {boolean} remove
 */ function skip_nodes() {
    let remove = arguments.length > 0 && arguments[0] !== void 0 ? arguments[0] : true;
    var depth = 0;
    var node = hydrate_node;
    while(true){
        if (node.nodeType === _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .COMMENT_NODE */.dz) {
            var data = /** @type {Comment} */ node.data;
            if (data === _constants_js__WEBPACK_IMPORTED_MODULE_1__/* .HYDRATION_END */.Lc) {
                if (depth === 0) return node;
                depth -= 1;
            } else if (data === _constants_js__WEBPACK_IMPORTED_MODULE_1__/* .HYDRATION_START */.CD || data === _constants_js__WEBPACK_IMPORTED_MODULE_1__/* .HYDRATION_START_ELSE */.qn) {
                depth += 1;
            }
        }
        var next = /** @type {TemplateNode} */ (0,_operations_js__WEBPACK_IMPORTED_MODULE_2__/* .get_next_sibling */.M$)(node);
        if (remove) node.remove();
        node = next;
    }
}
/**
 *
 * @param {TemplateNode} node
 */ function read_hydration_instruction(node) {
    if (!node || node.nodeType !== _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .COMMENT_NODE */.dz) {
        _warnings_js__WEBPACK_IMPORTED_MODULE_3__/* .hydration_mismatch */.eZ();
        throw _constants_js__WEBPACK_IMPORTED_MODULE_1__/* .HYDRATION_ERROR */.kD;
    }
    return /** @type {Comment} */ node.data;
}


}),
518: (function (__unused_webpack_module, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  Ey: () => (init_operations),
  Lo: () => (is_firefox),
  M$: () => (get_next_sibling),
  MC: () => (clear_text_content),
  Pb: () => (create_text),
  Zj: () => (get_first_child),
  eL: () => (should_defer_append),
  es: () => (first_child),
  hg: () => (sibling),
  jf: () => (child)
});
/* ESM import */var _hydration_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(452);
/* ESM import */var esm_env__WEBPACK_IMPORTED_MODULE_6__ = __webpack_require__(832);
/* ESM import */var _dev_equality_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(301);
/* ESM import */var _shared_utils_js__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(986);
/* ESM import */var _runtime_js__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(513);
/* ESM import */var _flags_index_js__WEBPACK_IMPORTED_MODULE_7__ = __webpack_require__(817);
/* ESM import */var _client_constants__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(924);
/* ESM import */var _reactivity_batch_js__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(410);
/** @import { Effect, TemplateNode } from '#client' */ 







// export these for reference in the compiled code, making global name deduplication unnecessary
/** @type {Window} */ var $window;
/** @type {Document} */ var $document;
/** @type {boolean} */ var is_firefox;
/** @type {() => Node | null} */ var first_child_getter;
/** @type {() => Node | null} */ var next_sibling_getter;
/**
 * Initialize these lazily to avoid issues when using the runtime in a server context
 * where these globals are not available while avoiding a separate server entry point
 */ function init_operations() {
    if ($window !== undefined) {
        return;
    }
    $window = window;
    $document = document;
    is_firefox = /Firefox/.test(navigator.userAgent);
    var element_prototype = Element.prototype;
    var node_prototype = Node.prototype;
    var text_prototype = Text.prototype;
    // @ts-ignore
    first_child_getter = (0,_shared_utils_js__WEBPACK_IMPORTED_MODULE_2__/* .get_descriptor */.J8)(node_prototype, 'firstChild').get;
    // @ts-ignore
    next_sibling_getter = (0,_shared_utils_js__WEBPACK_IMPORTED_MODULE_2__/* .get_descriptor */.J8)(node_prototype, 'nextSibling').get;
    if ((0,_shared_utils_js__WEBPACK_IMPORTED_MODULE_2__/* .is_extensible */.ZZ)(element_prototype)) {
        // the following assignments improve perf of lookups on DOM nodes
        // @ts-expect-error
        element_prototype.__click = undefined;
        // @ts-expect-error
        element_prototype.__className = undefined;
        // @ts-expect-error
        element_prototype.__attributes = null;
        // @ts-expect-error
        element_prototype.__style = undefined;
        // @ts-expect-error
        element_prototype.__e = undefined;
    }
    if ((0,_shared_utils_js__WEBPACK_IMPORTED_MODULE_2__/* .is_extensible */.ZZ)(text_prototype)) {
        // @ts-expect-error
        text_prototype.__t = undefined;
    }
    if (esm_env__WEBPACK_IMPORTED_MODULE_6__/* ["default"] */.A) {
        // @ts-expect-error
        element_prototype.__svelte_meta = null;
        (0,_dev_equality_js__WEBPACK_IMPORTED_MODULE_1__/* .init_array_prototype_warnings */.Ej)();
    }
}
/**
 * @param {string} value
 * @returns {Text}
 */ function create_text() {
    let value = arguments.length > 0 && arguments[0] !== void 0 ? arguments[0] : '';
    return document.createTextNode(value);
}
/**
 * @template {Node} N
 * @param {N} node
 * @returns {Node | null}
 */ /*@__NO_SIDE_EFFECTS__*/ function get_first_child(node) {
    return first_child_getter.call(node);
}
/**
 * @template {Node} N
 * @param {N} node
 * @returns {Node | null}
 */ /*@__NO_SIDE_EFFECTS__*/ function get_next_sibling(node) {
    return next_sibling_getter.call(node);
}
/**
 * Don't mark this as side-effect-free, hydration needs to walk all nodes
 * @template {Node} N
 * @param {N} node
 * @param {boolean} is_text
 * @returns {Node | null}
 */ function child(node, is_text) {
    if (!_hydration_js__WEBPACK_IMPORTED_MODULE_0__/* .hydrating */.fE) {
        return get_first_child(node);
    }
    var child = /** @type {TemplateNode} */ get_first_child(_hydration_js__WEBPACK_IMPORTED_MODULE_0__/* .hydrate_node */.Xb);
    // Child can be null if we have an element with a single child, like `<p>{text}</p>`, where `text` is empty
    if (child === null) {
        child = _hydration_js__WEBPACK_IMPORTED_MODULE_0__/* .hydrate_node.appendChild */.Xb.appendChild(create_text());
    } else if (is_text && child.nodeType !== _client_constants__WEBPACK_IMPORTED_MODULE_4__/* .TEXT_NODE */.Nd) {
        var text = create_text();
        child === null || child === void 0 ? void 0 : child.before(text);
        (0,_hydration_js__WEBPACK_IMPORTED_MODULE_0__/* .set_hydrate_node */.W0)(text);
        return text;
    }
    (0,_hydration_js__WEBPACK_IMPORTED_MODULE_0__/* .set_hydrate_node */.W0)(child);
    return child;
}
/**
 * Don't mark this as side-effect-free, hydration needs to walk all nodes
 * @param {DocumentFragment | TemplateNode | TemplateNode[]} fragment
 * @param {boolean} [is_text]
 * @returns {Node | null}
 */ function first_child(fragment) {
    let is_text = arguments.length > 1 && arguments[1] !== void 0 ? arguments[1] : false;
    if (!_hydration_js__WEBPACK_IMPORTED_MODULE_0__/* .hydrating */.fE) {
        // when not hydrating, `fragment` is a `DocumentFragment` (the result of calling `open_frag`)
        var first = /** @type {DocumentFragment} */ get_first_child(/** @type {Node} */ fragment);
        // TODO prevent user comments with the empty string when preserveComments is true
        if (first instanceof Comment && first.data === '') return get_next_sibling(first);
        return first;
    }
    // if an {expression} is empty during SSR, there might be no
    // text node to hydrate — we must therefore create one
    if (is_text && (_hydration_js__WEBPACK_IMPORTED_MODULE_0__/* .hydrate_node */.Xb === null || _hydration_js__WEBPACK_IMPORTED_MODULE_0__/* .hydrate_node */.Xb === void 0 ? void 0 : _hydration_js__WEBPACK_IMPORTED_MODULE_0__/* .hydrate_node.nodeType */.Xb.nodeType) !== _client_constants__WEBPACK_IMPORTED_MODULE_4__/* .TEXT_NODE */.Nd) {
        var text = create_text();
        _hydration_js__WEBPACK_IMPORTED_MODULE_0__/* .hydrate_node */.Xb === null || _hydration_js__WEBPACK_IMPORTED_MODULE_0__/* .hydrate_node */.Xb === void 0 ? void 0 : _hydration_js__WEBPACK_IMPORTED_MODULE_0__/* .hydrate_node.before */.Xb.before(text);
        (0,_hydration_js__WEBPACK_IMPORTED_MODULE_0__/* .set_hydrate_node */.W0)(text);
        return text;
    }
    return _hydration_js__WEBPACK_IMPORTED_MODULE_0__/* .hydrate_node */.Xb;
}
/**
 * Don't mark this as side-effect-free, hydration needs to walk all nodes
 * @param {TemplateNode} node
 * @param {number} count
 * @param {boolean} is_text
 * @returns {Node | null}
 */ function sibling(node) {
    let count = arguments.length > 1 && arguments[1] !== void 0 ? arguments[1] : 1, is_text = arguments.length > 2 && arguments[2] !== void 0 ? arguments[2] : false;
    let next_sibling = _hydration_js__WEBPACK_IMPORTED_MODULE_0__/* .hydrating */.fE ? _hydration_js__WEBPACK_IMPORTED_MODULE_0__/* .hydrate_node */.Xb : node;
    var last_sibling;
    while(count--){
        last_sibling = next_sibling;
        next_sibling = /** @type {TemplateNode} */ get_next_sibling(next_sibling);
    }
    if (!_hydration_js__WEBPACK_IMPORTED_MODULE_0__/* .hydrating */.fE) {
        return next_sibling;
    }
    // if a sibling {expression} is empty during SSR, there might be no
    // text node to hydrate — we must therefore create one
    if (is_text && (next_sibling === null || next_sibling === void 0 ? void 0 : next_sibling.nodeType) !== _client_constants__WEBPACK_IMPORTED_MODULE_4__/* .TEXT_NODE */.Nd) {
        var text = create_text();
        // If the next sibling is `null` and we're handling text then it's because
        // the SSR content was empty for the text, so we need to generate a new text
        // node and insert it after the last sibling
        if (next_sibling === null) {
            last_sibling === null || last_sibling === void 0 ? void 0 : last_sibling.after(text);
        } else {
            next_sibling.before(text);
        }
        (0,_hydration_js__WEBPACK_IMPORTED_MODULE_0__/* .set_hydrate_node */.W0)(text);
        return text;
    }
    (0,_hydration_js__WEBPACK_IMPORTED_MODULE_0__/* .set_hydrate_node */.W0)(next_sibling);
    return /** @type {TemplateNode} */ next_sibling;
}
/**
 * @template {Node} N
 * @param {N} node
 * @returns {void}
 */ function clear_text_content(node) {
    node.textContent = '';
}
/**
 * Returns `true` if we're updating the current block, for example `condition` in
 * an `{#if condition}` block just changed. In this case, the branch should be
 * appended (or removed) at the same time as other updates within the
 * current `<svelte:boundary>`
 */ function should_defer_append() {
    if (!_flags_index_js__WEBPACK_IMPORTED_MODULE_7__/* .async_mode_flag */.I0) return false;
    if (_reactivity_batch_js__WEBPACK_IMPORTED_MODULE_5__/* .eager_block_effects */.es !== null) return false;
    var flags = /** @type {Effect} */ _runtime_js__WEBPACK_IMPORTED_MODULE_3__/* .active_effect.f */.Fg.f;
    return (flags & _client_constants__WEBPACK_IMPORTED_MODULE_4__/* .EFFECT_RAN */.wi) !== 0;
}
/**
 *
 * @param {string} tag
 * @param {string} [namespace]
 * @param {string} [is]
 * @returns
 */ function create_element(tag, namespace, is) {
    let options = is ? {
        is
    } : undefined;
    if (namespace) {
        return document.createElementNS(namespace, tag, options);
    }
    return document.createElement(tag, options);
}
function create_fragment() {
    return document.createDocumentFragment();
}
/**
 * @param {string} data
 * @returns
 */ function create_comment() {
    let data = arguments.length > 0 && arguments[0] !== void 0 ? arguments[0] : '';
    return document.createComment(data);
}
/**
 * @param {Element} element
 * @param {string} key
 * @param {string} value
 * @returns
 */ function set_attribute(element, key) {
    let value = arguments.length > 2 && arguments[2] !== void 0 ? arguments[2] : '';
    if (key.startsWith('xlink:')) {
        element.setAttributeNS('http://www.w3.org/1999/xlink', key, value);
        return;
    }
    return element.setAttribute(key, value);
}


}),
642: (function (__unused_webpack_module, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  L: () => (create_fragment_from_html)
});
/** @param {string} html */ function create_fragment_from_html(html) {
    var elem = document.createElement('template');
    elem.innerHTML = html.replaceAll('<!>', '<!---->'); // XHTML compliance
    return elem.content;
}


}),
593: (function (__unused_webpack_module, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  $r: () => (queue_micro_task),
  KJ: () => (has_pending_tasks),
  eo: () => (flush_tasks)
});
/* ESM import */var _shared_utils_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(986);
/* ESM import */var _reactivity_batch_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(410);


// Fallback for when requestIdleCallback is not available
const request_idle_callback = (/* unused pure expression or super */ null && (typeof requestIdleCallback === 'undefined' ? (/** @type {() => void} */ cb)=>setTimeout(cb, 1) : requestIdleCallback));
/** @type {Array<() => void>} */ let micro_tasks = [];
/** @type {Array<() => void>} */ let idle_tasks = [];
function run_micro_tasks() {
    var tasks = micro_tasks;
    micro_tasks = [];
    (0,_shared_utils_js__WEBPACK_IMPORTED_MODULE_0__/* .run_all */.oO)(tasks);
}
function run_idle_tasks() {
    var tasks = idle_tasks;
    idle_tasks = [];
    (0,_shared_utils_js__WEBPACK_IMPORTED_MODULE_0__/* .run_all */.oO)(tasks);
}
function has_pending_tasks() {
    return micro_tasks.length > 0 || idle_tasks.length > 0;
}
/**
 * @param {() => void} fn
 */ function queue_micro_task(fn) {
    if (micro_tasks.length === 0 && !_reactivity_batch_js__WEBPACK_IMPORTED_MODULE_1__/* .is_flushing_sync */.OH) {
        var tasks = micro_tasks;
        queueMicrotask(()=>{
            // If this is false, a flushSync happened in the meantime. Do _not_ run new scheduled microtasks in that case
            // as the ordering of microtasks would be broken at that point - consider this case:
            // - queue_micro_task schedules microtask A to flush task X
            // - synchronously after, flushSync runs, processing task X
            // - synchronously after, some other microtask B is scheduled, but not through queue_micro_task but for example a Promise.resolve() in user code
            // - synchronously after, queue_micro_task schedules microtask C to flush task Y
            // - one tick later, microtask A now resolves, flushing task Y before microtask B, which is incorrect
            // This if check prevents that race condition (that realistically will only happen in tests)
            if (tasks === micro_tasks) run_micro_tasks();
        });
    }
    micro_tasks.push(fn);
}
/**
 * @param {() => void} fn
 */ function queue_idle_task(fn) {
    if (idle_tasks.length === 0) {
        request_idle_callback(run_idle_tasks);
    }
    idle_tasks.push(fn);
}
/**
 * Synchronously run any queued tasks.
 */ function flush_tasks() {
    if (micro_tasks.length > 0) {
        run_micro_tasks();
    }
    if (idle_tasks.length > 0) {
        run_idle_tasks();
    }
}


}),
782: (function (__unused_webpack_module, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  BC: () => (append),
  Im: () => (comment),
  mX: () => (assign_nodes),
  vU: () => (from_html)
});
/* ESM import */var _hydration_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(452);
/* ESM import */var _operations_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(518);
/* ESM import */var _reconciler_js__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(642);
/* ESM import */var _runtime_js__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(513);
/* ESM import */var _constants_js__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(178);
/* ESM import */var _client_constants__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(924);
/** @import { Effect, TemplateNode } from '#client' */ /** @import { TemplateStructure } from './types' */ 





/**
 * @param {TemplateNode} start
 * @param {TemplateNode | null} end
 */ function assign_nodes(start, end) {
    var effect = /** @type {Effect} */ _runtime_js__WEBPACK_IMPORTED_MODULE_2__/* .active_effect */.Fg;
    if (effect.nodes_start === null) {
        effect.nodes_start = start;
        effect.nodes_end = end;
    }
}
/**
 * @param {string} content
 * @param {number} flags
 * @returns {() => Node | Node[]}
 */ /*#__NO_SIDE_EFFECTS__*/ function from_html(content, flags) {
    var is_fragment = (flags & _constants_js__WEBPACK_IMPORTED_MODULE_3__/* .TEMPLATE_FRAGMENT */.Ax) !== 0;
    var use_import_node = (flags & _constants_js__WEBPACK_IMPORTED_MODULE_3__/* .TEMPLATE_USE_IMPORT_NODE */.iX) !== 0;
    /** @type {Node} */ var node;
    /**
	 * Whether or not the first item is a text/element node. If not, we need to
	 * create an additional comment node to act as `effect.nodes.start`
	 */ var has_start = !content.startsWith('<!>');
    return ()=>{
        if (_hydration_js__WEBPACK_IMPORTED_MODULE_0__/* .hydrating */.fE) {
            assign_nodes(_hydration_js__WEBPACK_IMPORTED_MODULE_0__/* .hydrate_node */.Xb, null);
            return _hydration_js__WEBPACK_IMPORTED_MODULE_0__/* .hydrate_node */.Xb;
        }
        if (node === undefined) {
            node = (0,_reconciler_js__WEBPACK_IMPORTED_MODULE_5__/* .create_fragment_from_html */.L)(has_start ? content : '<!>' + content);
            if (!is_fragment) node = /** @type {Node} */ (0,_operations_js__WEBPACK_IMPORTED_MODULE_1__/* .get_first_child */.Zj)(node);
        }
        var clone = /** @type {TemplateNode} */ use_import_node || _operations_js__WEBPACK_IMPORTED_MODULE_1__/* .is_firefox */.Lo ? document.importNode(node, true) : node.cloneNode(true);
        if (is_fragment) {
            var start = /** @type {TemplateNode} */ (0,_operations_js__WEBPACK_IMPORTED_MODULE_1__/* .get_first_child */.Zj)(clone);
            var end = /** @type {TemplateNode} */ clone.lastChild;
            assign_nodes(start, end);
        } else {
            assign_nodes(clone, clone);
        }
        return clone;
    };
}
/**
 * @param {string} content
 * @param {number} flags
 * @param {'svg' | 'math'} ns
 * @returns {() => Node | Node[]}
 */ /*#__NO_SIDE_EFFECTS__*/ function from_namespace(content, flags) {
    let ns = arguments.length > 2 && arguments[2] !== void 0 ? arguments[2] : 'svg';
    /**
	 * Whether or not the first item is a text/element node. If not, we need to
	 * create an additional comment node to act as `effect.nodes.start`
	 */ var has_start = !content.startsWith('<!>');
    var is_fragment = (flags & TEMPLATE_FRAGMENT) !== 0;
    var wrapped = `<${ns}>${has_start ? content : '<!>' + content}</${ns}>`;
    /** @type {Element | DocumentFragment} */ var node;
    return ()=>{
        if (hydrating) {
            assign_nodes(hydrate_node, null);
            return hydrate_node;
        }
        if (!node) {
            var fragment = /** @type {DocumentFragment} */ create_fragment_from_html(wrapped);
            var root = /** @type {Element} */ get_first_child(fragment);
            if (is_fragment) {
                node = document.createDocumentFragment();
                while(get_first_child(root)){
                    node.appendChild(/** @type {Node} */ get_first_child(root));
                }
            } else {
                node = /** @type {Element} */ get_first_child(root);
            }
        }
        var clone = /** @type {TemplateNode} */ node.cloneNode(true);
        if (is_fragment) {
            var start = /** @type {TemplateNode} */ get_first_child(clone);
            var end = /** @type {TemplateNode} */ clone.lastChild;
            assign_nodes(start, end);
        } else {
            assign_nodes(clone, clone);
        }
        return clone;
    };
}
/**
 * @param {string} content
 * @param {number} flags
 */ /*#__NO_SIDE_EFFECTS__*/ function from_svg(content, flags) {
    return from_namespace(content, flags, 'svg');
}
/**
 * @param {string} content
 * @param {number} flags
 */ /*#__NO_SIDE_EFFECTS__*/ function from_mathml(content, flags) {
    return from_namespace(content, flags, 'math');
}
/**
 * @param {TemplateStructure[]} structure
 * @param {typeof NAMESPACE_SVG | typeof NAMESPACE_MATHML | undefined} [ns]
 */ function fragment_from_tree(structure, ns) {
    var fragment = create_fragment();
    for (var item of structure){
        if (typeof item === 'string') {
            fragment.append(create_text(item));
            continue;
        }
        // if `preserveComments === true`, comments are represented as `['// <data>']`
        if (item === undefined || item[0][0] === '/') {
            fragment.append(create_comment(item ? item[0].slice(3) : ''));
            continue;
        }
        const [name, attributes, ...children] = item;
        const namespace = name === 'svg' ? NAMESPACE_SVG : name === 'math' ? NAMESPACE_MATHML : ns;
        var element = create_element(name, namespace, attributes === null || attributes === void 0 ? void 0 : attributes.is);
        for(var key in attributes){
            set_attribute(element, key, attributes[key]);
        }
        if (children.length > 0) {
            var target = element.tagName === 'TEMPLATE' ? /** @type {HTMLTemplateElement} */ element.content : element;
            target.append(fragment_from_tree(children, element.tagName === 'foreignObject' ? undefined : namespace));
        }
        fragment.append(element);
    }
    return fragment;
}
/**
 * @param {TemplateStructure[]} structure
 * @param {number} flags
 * @returns {() => Node | Node[]}
 */ /*#__NO_SIDE_EFFECTS__*/ function from_tree(structure, flags) {
    var is_fragment = (flags & TEMPLATE_FRAGMENT) !== 0;
    var use_import_node = (flags & TEMPLATE_USE_IMPORT_NODE) !== 0;
    /** @type {Node} */ var node;
    return ()=>{
        if (hydrating) {
            assign_nodes(hydrate_node, null);
            return hydrate_node;
        }
        if (node === undefined) {
            const ns = (flags & TEMPLATE_USE_SVG) !== 0 ? NAMESPACE_SVG : (flags & TEMPLATE_USE_MATHML) !== 0 ? NAMESPACE_MATHML : undefined;
            node = fragment_from_tree(structure, ns);
            if (!is_fragment) node = /** @type {Node} */ get_first_child(node);
        }
        var clone = /** @type {TemplateNode} */ use_import_node || is_firefox ? document.importNode(node, true) : node.cloneNode(true);
        if (is_fragment) {
            var start = /** @type {TemplateNode} */ get_first_child(clone);
            var end = /** @type {TemplateNode} */ clone.lastChild;
            assign_nodes(start, end);
        } else {
            assign_nodes(clone, clone);
        }
        return clone;
    };
}
/**
 * @param {() => Element | DocumentFragment} fn
 */ function with_script(fn) {
    return ()=>run_scripts(fn());
}
/**
 * Creating a document fragment from HTML that contains script tags will not execute
 * the scripts. We need to replace the script tags with new ones so that they are executed.
 * @param {Element | DocumentFragment} node
 * @returns {Node | Node[]}
 */ function run_scripts(node) {
    // scripts were SSR'd, in which case they will run
    if (hydrating) return node;
    const is_fragment = node.nodeType === DOCUMENT_FRAGMENT_NODE;
    const scripts = /** @type {HTMLElement} */ node.tagName === 'SCRIPT' ? [
        /** @type {HTMLScriptElement} */ node
    ] : node.querySelectorAll('script');
    const effect = /** @type {Effect} */ active_effect;
    for (const script of scripts){
        const clone = document.createElement('script');
        for (var attribute of script.attributes){
            clone.setAttribute(attribute.name, attribute.value);
        }
        clone.textContent = script.textContent;
        // The script has changed - if it's at the edges, the effect now points at dead nodes
        if (is_fragment ? node.firstChild === script : node === script) {
            effect.nodes_start = clone;
        }
        if (is_fragment ? node.lastChild === script : node === script) {
            effect.nodes_end = clone;
        }
        script.replaceWith(clone);
    }
    return node;
}
/**
 * Don't mark this as side-effect-free, hydration needs to walk all nodes
 * @param {any} value
 */ function text() {
    let value = arguments.length > 0 && arguments[0] !== void 0 ? arguments[0] : '';
    if (!hydrating) {
        var t = create_text(value + '');
        assign_nodes(t, t);
        return t;
    }
    var node = hydrate_node;
    if (node.nodeType !== TEXT_NODE) {
        // if an {expression} is empty during SSR, we need to insert an empty text node
        node.before(node = create_text());
        set_hydrate_node(node);
    }
    assign_nodes(node, node);
    return node;
}
/**
 * @returns {TemplateNode | DocumentFragment}
 */ function comment() {
    // we're not delegating to `template` here for performance reasons
    if (_hydration_js__WEBPACK_IMPORTED_MODULE_0__/* .hydrating */.fE) {
        assign_nodes(_hydration_js__WEBPACK_IMPORTED_MODULE_0__/* .hydrate_node */.Xb, null);
        return _hydration_js__WEBPACK_IMPORTED_MODULE_0__/* .hydrate_node */.Xb;
    }
    var frag = document.createDocumentFragment();
    var start = document.createComment('');
    var anchor = (0,_operations_js__WEBPACK_IMPORTED_MODULE_1__/* .create_text */.Pb)();
    frag.append(start, anchor);
    assign_nodes(start, anchor);
    return frag;
}
/**
 * Assign the created (or in hydration mode, traversed) dom elements to the current block
 * and insert the elements into the dom (in client mode).
 * @param {Text | Comment | Element} anchor
 * @param {DocumentFragment | Element} dom
 */ function append(anchor, dom) {
    if (_hydration_js__WEBPACK_IMPORTED_MODULE_0__/* .hydrating */.fE) {
        /** @type {Effect} */ _runtime_js__WEBPACK_IMPORTED_MODULE_2__/* .active_effect.nodes_end */.Fg.nodes_end = _hydration_js__WEBPACK_IMPORTED_MODULE_0__/* .hydrate_node */.Xb;
        (0,_hydration_js__WEBPACK_IMPORTED_MODULE_0__/* .hydrate_next */.E$)();
        return;
    }
    if (anchor === null) {
        // edge case — void `<svelte:element>` with content
        return;
    }
    anchor.before(/** @type {Node} */ dom);
}
/**
 * Create (or hydrate) an unique UID for the component instance.
 */ function props_id() {
    var _hydrate_node_textContent;
    var // @ts-expect-error This way we ensure the id is unique even across Svelte runtimes
    _ref, _window;
    if (hydrating && hydrate_node && hydrate_node.nodeType === COMMENT_NODE && ((_hydrate_node_textContent = hydrate_node.textContent) === null || _hydrate_node_textContent === void 0 ? void 0 : _hydrate_node_textContent.startsWith(`#`))) {
        const id = hydrate_node.textContent.substring(1);
        hydrate_next();
        return id;
    }
    (_ref = (_window = window).__svelte ?? (_window.__svelte = {})).uid ?? (_ref.uid = 1);
    // @ts-expect-error
    return `c${window.__svelte.uid++}`;
}


}),
621: (function (__unused_webpack_module, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  i: () => (handle_error),
  n: () => (invoke_error_boundary)
});
/* ESM import */var esm_env__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(832);
/* ESM import */var _constants_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(178);
/* ESM import */var _dom_operations_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(518);
/* ESM import */var _constants_js__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(924);
/* ESM import */var _shared_utils_js__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(986);
/* ESM import */var _runtime_js__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(513);
/** @import { Derived, Effect } from '#client' */ /** @import { Boundary } from './dom/blocks/boundary.js' */ 





const adjustments = new WeakMap();
/**
 * @param {unknown} error
 */ function handle_error(error) {
    var effect = _runtime_js__WEBPACK_IMPORTED_MODULE_4__/* .active_effect */.Fg;
    // for unowned deriveds, don't throw until we read the value
    if (effect === null) {
        /** @type {Derived} */ _runtime_js__WEBPACK_IMPORTED_MODULE_4__/* .active_reaction.f */.hp.f |= _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .ERROR_VALUE */.dH;
        return error;
    }
    if (esm_env__WEBPACK_IMPORTED_MODULE_5__/* ["default"] */.A && error instanceof Error && !adjustments.has(error)) {
        adjustments.set(error, get_adjustments(error, effect));
    }
    if ((effect.f & _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .EFFECT_RAN */.wi) === 0) {
        // if the error occurred while creating this subtree, we let it
        // bubble up until it hits a boundary that can handle it
        if ((effect.f & _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .BOUNDARY_EFFECT */.bp) === 0) {
            if (!effect.parent && error instanceof Error) {
                apply_adjustments(error);
            }
            throw error;
        }
        /** @type {Boundary} */ effect.b.error(error);
    } else {
        // otherwise we bubble up the effect tree ourselves
        invoke_error_boundary(error, effect);
    }
}
/**
 * @param {unknown} error
 * @param {Effect | null} effect
 */ function invoke_error_boundary(error, effect) {
    while(effect !== null){
        if ((effect.f & _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .BOUNDARY_EFFECT */.bp) !== 0) {
            try {
                /** @type {Boundary} */ effect.b.error(error);
                return;
            } catch (e) {
                error = e;
            }
        }
        effect = effect.parent;
    }
    if (error instanceof Error) {
        apply_adjustments(error);
    }
    throw error;
}
/**
 * Add useful information to the error message/stack in development
 * @param {Error} error
 * @param {Effect} effect
 */ function get_adjustments(error, effect) {
    var _effect_fn, _error_stack;
    const message_descriptor = (0,_shared_utils_js__WEBPACK_IMPORTED_MODULE_3__/* .get_descriptor */.J8)(error, 'message');
    // if the message was already changed and it's not configurable we can't change it
    // or it will throw a different error swallowing the original error
    if (message_descriptor && !message_descriptor.configurable) return;
    var indent = _dom_operations_js__WEBPACK_IMPORTED_MODULE_1__/* .is_firefox */.Lo ? '  ' : '\t';
    var component_stack = `\n${indent}in ${((_effect_fn = effect.fn) === null || _effect_fn === void 0 ? void 0 : _effect_fn.name) || '<unknown>'}`;
    var context = effect.ctx;
    while(context !== null){
        var _context_function;
        component_stack += `\n${indent}in ${(_context_function = context.function) === null || _context_function === void 0 ? void 0 : _context_function[_constants_js__WEBPACK_IMPORTED_MODULE_0__/* .FILENAME */.Uh].split('/').pop()}`;
        context = context.p;
    }
    return {
        message: error.message + `\n${component_stack}\n`,
        stack: (_error_stack = error.stack) === null || _error_stack === void 0 ? void 0 : _error_stack.split('\n').filter((line)=>!line.includes('svelte/src/internal')).join('\n')
    };
}
/**
 * @param {Error} error
 */ function apply_adjustments(error) {
    const adjusted = adjustments.get(error);
    if (adjusted) {
        (0,_shared_utils_js__WEBPACK_IMPORTED_MODULE_3__/* .define_property */.Qu)(error, 'message', {
            value: adjusted.message
        });
        (0,_shared_utils_js__WEBPACK_IMPORTED_MODULE_3__/* .define_property */.Qu)(error, 'stack', {
            value: adjusted.stack
        });
    }
}


}),
626: (function (__unused_webpack_module, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  BT: () => (effect_in_teardown),
  Cl: () => (effect_update_depth_exceeded),
  JJ: () => (svelte_boundary_reset_onerror),
  Uw: () => (state_descriptors_fixed),
  Vv: () => (hydration_failed),
  WR: () => (invalid_snippet),
  YY: () => (state_prototype_fixed),
  aQ: () => (async_derived_orphan),
  cN: () => (derived_references_self),
  fW: () => (flush_sync_in_effect),
  fi: () => (effect_in_unowned_derived),
  js: () => (props_rest_readonly),
  rZ: () => (state_unsafe_mutation),
  tB: () => (effect_orphan),
  xU: () => (rune_outside_svelte)
});
/* ESM import */var esm_env__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(832);
/* This file is generated by scripts/process-messages/index.js. Do not edit! */ 

/**
 * Cannot create a `$derived(...)` with an `await` expression outside of an effect tree
 * @returns {never}
 */ function async_derived_orphan() {
    if (esm_env__WEBPACK_IMPORTED_MODULE_0__/* ["default"] */.A) {
        const error = new Error(`async_derived_orphan\nCannot create a \`$derived(...)\` with an \`await\` expression outside of an effect tree\nhttps://svelte.dev/e/async_derived_orphan`);
        error.name = 'Svelte error';
        throw error;
    } else {
        throw new Error(`https://svelte.dev/e/async_derived_orphan`);
    }
}
/**
 * Using `bind:value` together with a checkbox input is not allowed. Use `bind:checked` instead
 * @returns {never}
 */ function bind_invalid_checkbox_value() {
    if (DEV) {
        const error = new Error(`bind_invalid_checkbox_value\nUsing \`bind:value\` together with a checkbox input is not allowed. Use \`bind:checked\` instead\nhttps://svelte.dev/e/bind_invalid_checkbox_value`);
        error.name = 'Svelte error';
        throw error;
    } else {
        throw new Error(`https://svelte.dev/e/bind_invalid_checkbox_value`);
    }
}
/**
 * Component %component% has an export named `%key%` that a consumer component is trying to access using `bind:%key%`, which is disallowed. Instead, use `bind:this` (e.g. `<%name% bind:this={component} />`) and then access the property on the bound component instance (e.g. `component.%key%`)
 * @param {string} component
 * @param {string} key
 * @param {string} name
 * @returns {never}
 */ function bind_invalid_export(component, key, name) {
    if (DEV) {
        const error = new Error(`bind_invalid_export\nComponent ${component} has an export named \`${key}\` that a consumer component is trying to access using \`bind:${key}\`, which is disallowed. Instead, use \`bind:this\` (e.g. \`<${name} bind:this={component} />\`) and then access the property on the bound component instance (e.g. \`component.${key}\`)\nhttps://svelte.dev/e/bind_invalid_export`);
        error.name = 'Svelte error';
        throw error;
    } else {
        throw new Error(`https://svelte.dev/e/bind_invalid_export`);
    }
}
/**
 * A component is attempting to bind to a non-bindable property `%key%` belonging to %component% (i.e. `<%name% bind:%key%={...}>`). To mark a property as bindable: `let { %key% = $bindable() } = $props()`
 * @param {string} key
 * @param {string} component
 * @param {string} name
 * @returns {never}
 */ function bind_not_bindable(key, component, name) {
    if (DEV) {
        const error = new Error(`bind_not_bindable\nA component is attempting to bind to a non-bindable property \`${key}\` belonging to ${component} (i.e. \`<${name} bind:${key}={...}>\`). To mark a property as bindable: \`let { ${key} = $bindable() } = $props()\`\nhttps://svelte.dev/e/bind_not_bindable`);
        error.name = 'Svelte error';
        throw error;
    } else {
        throw new Error(`https://svelte.dev/e/bind_not_bindable`);
    }
}
/**
 * Calling `%method%` on a component instance (of %component%) is no longer valid in Svelte 5
 * @param {string} method
 * @param {string} component
 * @returns {never}
 */ function component_api_changed(method, component) {
    if (DEV) {
        const error = new Error(`component_api_changed\nCalling \`${method}\` on a component instance (of ${component}) is no longer valid in Svelte 5\nhttps://svelte.dev/e/component_api_changed`);
        error.name = 'Svelte error';
        throw error;
    } else {
        throw new Error(`https://svelte.dev/e/component_api_changed`);
    }
}
/**
 * Attempted to instantiate %component% with `new %name%`, which is no longer valid in Svelte 5. If this component is not under your control, set the `compatibility.componentApi` compiler option to `4` to keep it working.
 * @param {string} component
 * @param {string} name
 * @returns {never}
 */ function component_api_invalid_new(component, name) {
    if (DEV) {
        const error = new Error(`component_api_invalid_new\nAttempted to instantiate ${component} with \`new ${name}\`, which is no longer valid in Svelte 5. If this component is not under your control, set the \`compatibility.componentApi\` compiler option to \`4\` to keep it working.\nhttps://svelte.dev/e/component_api_invalid_new`);
        error.name = 'Svelte error';
        throw error;
    } else {
        throw new Error(`https://svelte.dev/e/component_api_invalid_new`);
    }
}
/**
 * A derived value cannot reference itself recursively
 * @returns {never}
 */ function derived_references_self() {
    if (esm_env__WEBPACK_IMPORTED_MODULE_0__/* ["default"] */.A) {
        const error = new Error(`derived_references_self\nA derived value cannot reference itself recursively\nhttps://svelte.dev/e/derived_references_self`);
        error.name = 'Svelte error';
        throw error;
    } else {
        throw new Error(`https://svelte.dev/e/derived_references_self`);
    }
}
/**
 * Keyed each block has duplicate key `%value%` at indexes %a% and %b%
 * @param {string} a
 * @param {string} b
 * @param {string | undefined | null} [value]
 * @returns {never}
 */ function each_key_duplicate(a, b, value) {
    if (DEV) {
        const error = new Error(`each_key_duplicate\n${value ? `Keyed each block has duplicate key \`${value}\` at indexes ${a} and ${b}` : `Keyed each block has duplicate key at indexes ${a} and ${b}`}\nhttps://svelte.dev/e/each_key_duplicate`);
        error.name = 'Svelte error';
        throw error;
    } else {
        throw new Error(`https://svelte.dev/e/each_key_duplicate`);
    }
}
/**
 * `%rune%` cannot be used inside an effect cleanup function
 * @param {string} rune
 * @returns {never}
 */ function effect_in_teardown(rune) {
    if (esm_env__WEBPACK_IMPORTED_MODULE_0__/* ["default"] */.A) {
        const error = new Error(`effect_in_teardown\n\`${rune}\` cannot be used inside an effect cleanup function\nhttps://svelte.dev/e/effect_in_teardown`);
        error.name = 'Svelte error';
        throw error;
    } else {
        throw new Error(`https://svelte.dev/e/effect_in_teardown`);
    }
}
/**
 * Effect cannot be created inside a `$derived` value that was not itself created inside an effect
 * @returns {never}
 */ function effect_in_unowned_derived() {
    if (esm_env__WEBPACK_IMPORTED_MODULE_0__/* ["default"] */.A) {
        const error = new Error(`effect_in_unowned_derived\nEffect cannot be created inside a \`$derived\` value that was not itself created inside an effect\nhttps://svelte.dev/e/effect_in_unowned_derived`);
        error.name = 'Svelte error';
        throw error;
    } else {
        throw new Error(`https://svelte.dev/e/effect_in_unowned_derived`);
    }
}
/**
 * `%rune%` can only be used inside an effect (e.g. during component initialisation)
 * @param {string} rune
 * @returns {never}
 */ function effect_orphan(rune) {
    if (esm_env__WEBPACK_IMPORTED_MODULE_0__/* ["default"] */.A) {
        const error = new Error(`effect_orphan\n\`${rune}\` can only be used inside an effect (e.g. during component initialisation)\nhttps://svelte.dev/e/effect_orphan`);
        error.name = 'Svelte error';
        throw error;
    } else {
        throw new Error(`https://svelte.dev/e/effect_orphan`);
    }
}
/**
 * `$effect.pending()` can only be called inside an effect or derived
 * @returns {never}
 */ function effect_pending_outside_reaction() {
    if (DEV) {
        const error = new Error(`effect_pending_outside_reaction\n\`$effect.pending()\` can only be called inside an effect or derived\nhttps://svelte.dev/e/effect_pending_outside_reaction`);
        error.name = 'Svelte error';
        throw error;
    } else {
        throw new Error(`https://svelte.dev/e/effect_pending_outside_reaction`);
    }
}
/**
 * Maximum update depth exceeded. This typically indicates that an effect reads and writes the same piece of state
 * @returns {never}
 */ function effect_update_depth_exceeded() {
    if (esm_env__WEBPACK_IMPORTED_MODULE_0__/* ["default"] */.A) {
        const error = new Error(`effect_update_depth_exceeded\nMaximum update depth exceeded. This typically indicates that an effect reads and writes the same piece of state\nhttps://svelte.dev/e/effect_update_depth_exceeded`);
        error.name = 'Svelte error';
        throw error;
    } else {
        throw new Error(`https://svelte.dev/e/effect_update_depth_exceeded`);
    }
}
/**
 * Cannot use `flushSync` inside an effect
 * @returns {never}
 */ function flush_sync_in_effect() {
    if (esm_env__WEBPACK_IMPORTED_MODULE_0__/* ["default"] */.A) {
        const error = new Error(`flush_sync_in_effect\nCannot use \`flushSync\` inside an effect\nhttps://svelte.dev/e/flush_sync_in_effect`);
        error.name = 'Svelte error';
        throw error;
    } else {
        throw new Error(`https://svelte.dev/e/flush_sync_in_effect`);
    }
}
/**
 * `getAbortSignal()` can only be called inside an effect or derived
 * @returns {never}
 */ function get_abort_signal_outside_reaction() {
    if (DEV) {
        const error = new Error(`get_abort_signal_outside_reaction\n\`getAbortSignal()\` can only be called inside an effect or derived\nhttps://svelte.dev/e/get_abort_signal_outside_reaction`);
        error.name = 'Svelte error';
        throw error;
    } else {
        throw new Error(`https://svelte.dev/e/get_abort_signal_outside_reaction`);
    }
}
/**
 * Failed to hydrate the application
 * @returns {never}
 */ function hydration_failed() {
    if (esm_env__WEBPACK_IMPORTED_MODULE_0__/* ["default"] */.A) {
        const error = new Error(`hydration_failed\nFailed to hydrate the application\nhttps://svelte.dev/e/hydration_failed`);
        error.name = 'Svelte error';
        throw error;
    } else {
        throw new Error(`https://svelte.dev/e/hydration_failed`);
    }
}
/**
 * Could not `{@render}` snippet due to the expression being `null` or `undefined`. Consider using optional chaining `{@render snippet?.()}`
 * @returns {never}
 */ function invalid_snippet() {
    if (esm_env__WEBPACK_IMPORTED_MODULE_0__/* ["default"] */.A) {
        const error = new Error(`invalid_snippet\nCould not \`{@render}\` snippet due to the expression being \`null\` or \`undefined\`. Consider using optional chaining \`{@render snippet?.()}\`\nhttps://svelte.dev/e/invalid_snippet`);
        error.name = 'Svelte error';
        throw error;
    } else {
        throw new Error(`https://svelte.dev/e/invalid_snippet`);
    }
}
/**
 * `%name%(...)` cannot be used in runes mode
 * @param {string} name
 * @returns {never}
 */ function lifecycle_legacy_only(name) {
    if (DEV) {
        const error = new Error(`lifecycle_legacy_only\n\`${name}(...)\` cannot be used in runes mode\nhttps://svelte.dev/e/lifecycle_legacy_only`);
        error.name = 'Svelte error';
        throw error;
    } else {
        throw new Error(`https://svelte.dev/e/lifecycle_legacy_only`);
    }
}
/**
 * Cannot do `bind:%key%={undefined}` when `%key%` has a fallback value
 * @param {string} key
 * @returns {never}
 */ function props_invalid_value(key) {
    if (DEV) {
        const error = new Error(`props_invalid_value\nCannot do \`bind:${key}={undefined}\` when \`${key}\` has a fallback value\nhttps://svelte.dev/e/props_invalid_value`);
        error.name = 'Svelte error';
        throw error;
    } else {
        throw new Error(`https://svelte.dev/e/props_invalid_value`);
    }
}
/**
 * Rest element properties of `$props()` such as `%property%` are readonly
 * @param {string} property
 * @returns {never}
 */ function props_rest_readonly(property) {
    if (esm_env__WEBPACK_IMPORTED_MODULE_0__/* ["default"] */.A) {
        const error = new Error(`props_rest_readonly\nRest element properties of \`$props()\` such as \`${property}\` are readonly\nhttps://svelte.dev/e/props_rest_readonly`);
        error.name = 'Svelte error';
        throw error;
    } else {
        throw new Error(`https://svelte.dev/e/props_rest_readonly`);
    }
}
/**
 * The `%rune%` rune is only available inside `.svelte` and `.svelte.js/ts` files
 * @param {string} rune
 * @returns {never}
 */ function rune_outside_svelte(rune) {
    if (esm_env__WEBPACK_IMPORTED_MODULE_0__/* ["default"] */.A) {
        const error = new Error(`rune_outside_svelte\nThe \`${rune}\` rune is only available inside \`.svelte\` and \`.svelte.js/ts\` files\nhttps://svelte.dev/e/rune_outside_svelte`);
        error.name = 'Svelte error';
        throw error;
    } else {
        throw new Error(`https://svelte.dev/e/rune_outside_svelte`);
    }
}
/**
 * `setContext` must be called when a component first initializes, not in a subsequent effect or after an `await` expression
 * @returns {never}
 */ function set_context_after_init() {
    if (DEV) {
        const error = new Error(`set_context_after_init\n\`setContext\` must be called when a component first initializes, not in a subsequent effect or after an \`await\` expression\nhttps://svelte.dev/e/set_context_after_init`);
        error.name = 'Svelte error';
        throw error;
    } else {
        throw new Error(`https://svelte.dev/e/set_context_after_init`);
    }
}
/**
 * Property descriptors defined on `$state` objects must contain `value` and always be `enumerable`, `configurable` and `writable`.
 * @returns {never}
 */ function state_descriptors_fixed() {
    if (esm_env__WEBPACK_IMPORTED_MODULE_0__/* ["default"] */.A) {
        const error = new Error(`state_descriptors_fixed\nProperty descriptors defined on \`$state\` objects must contain \`value\` and always be \`enumerable\`, \`configurable\` and \`writable\`.\nhttps://svelte.dev/e/state_descriptors_fixed`);
        error.name = 'Svelte error';
        throw error;
    } else {
        throw new Error(`https://svelte.dev/e/state_descriptors_fixed`);
    }
}
/**
 * Cannot set prototype of `$state` object
 * @returns {never}
 */ function state_prototype_fixed() {
    if (esm_env__WEBPACK_IMPORTED_MODULE_0__/* ["default"] */.A) {
        const error = new Error(`state_prototype_fixed\nCannot set prototype of \`$state\` object\nhttps://svelte.dev/e/state_prototype_fixed`);
        error.name = 'Svelte error';
        throw error;
    } else {
        throw new Error(`https://svelte.dev/e/state_prototype_fixed`);
    }
}
/**
 * Updating state inside `$derived(...)`, `$inspect(...)` or a template expression is forbidden. If the value should not be reactive, declare it without `$state`
 * @returns {never}
 */ function state_unsafe_mutation() {
    if (esm_env__WEBPACK_IMPORTED_MODULE_0__/* ["default"] */.A) {
        const error = new Error(`state_unsafe_mutation\nUpdating state inside \`$derived(...)\`, \`$inspect(...)\` or a template expression is forbidden. If the value should not be reactive, declare it without \`$state\`\nhttps://svelte.dev/e/state_unsafe_mutation`);
        error.name = 'Svelte error';
        throw error;
    } else {
        throw new Error(`https://svelte.dev/e/state_unsafe_mutation`);
    }
}
/**
 * A `<svelte:boundary>` `reset` function cannot be called while an error is still being handled
 * @returns {never}
 */ function svelte_boundary_reset_onerror() {
    if (esm_env__WEBPACK_IMPORTED_MODULE_0__/* ["default"] */.A) {
        const error = new Error(`svelte_boundary_reset_onerror\nA \`<svelte:boundary>\` \`reset\` function cannot be called while an error is still being handled\nhttps://svelte.dev/e/svelte_boundary_reset_onerror`);
        error.name = 'Svelte error';
        throw error;
    } else {
        throw new Error(`https://svelte.dev/e/svelte_boundary_reset_onerror`);
    }
}


}),
750: (function (__unused_webpack_module, __webpack_exports__, __webpack_require__) {

// EXPORTS
__webpack_require__.d(__webpack_exports__, {
  wk1: () => (/* reexport */ reactivity_sources/* .state */.wk),
  "if": () => (/* reexport */ if_block),
  VCO: () => (/* reexport */ client_context/* .push */.VC),
  jax: () => (/* reexport */ render/* .set_text */.j),
  K2T: () => (/* reexport */ hydration/* .next */.K2),
  unG: () => (/* reexport */ deriveds/* .user_derived */.eO),
  JtY: () => (/* reexport */ runtime/* .get */.Jt),
  iRd: () => (/* reexport */ rest_props),
  zgK: () => (/* reexport */ reactivity_sources/* .mutable_source */.zg),
  hZp: () => (/* reexport */ reactivity_sources/* .set */.hZ),
  XId: () => (/* reexport */ actions_action),
  uYY: () => (/* reexport */ client_context/* .pop */.uY),
  jfp: () => (/* reexport */ operations/* .child */.jf),
  esp: () => (/* reexport */ operations/* .first_child */.es),
  cLc: () => (/* reexport */ hydration/* .reset */.cL),
  vNg: () => (/* reexport */ reactivity_effects/* .template_effect */.vN),
  Imx: () => (/* reexport */ template/* .comment */.Im),
  hg4: () => (/* reexport */ operations/* .sibling */.hg),
  MmH: () => (/* reexport */ elements_events/* .delegate */.Mm),
  vUu: () => (/* reexport */ template/* .from_html */.vU),
  qyt: () => (/* reexport */ html_html),
  UAl: () => (/* reexport */ snippet/* .snippet */.UA),
  BCw: () => (/* reexport */ template/* .append */.BC),
  MWq: () => (/* reexport */ reactivity_effects/* .user_effect */.MW)
});

// UNUSED EXPORTS: bind_window_size, deep_read_state, update_pre_prop, with_script, attr, pending, bind_select_value, css_props, sanitize_slots, update, set_style, update_pre, mutate, animation, bind_played, replay_events, bind_muted, active_effect, attach, from_mathml, tag_proxy, component, bind_online, legacy_pre_effect, store_set, each, text, FILENAME, event, props_id, inspect, validate_binding, legacy_api, setup_stores, cleanup_styles, deep_read, once, apply, bind_focused, legacy_rest_props, validate_void_dynamic_element, create_custom_element, track_reactivity_loss, bind_value, invoke_error_boundary, to_array, window, bind_paused, STYLE, bind_ready_state, from_svg, init, assign_or, key, async_body, effect_tracking, invalid_default_snippet, invalidate_inner_signals, head, store_mutate, hmr, stopPropagation, fallback, stopImmediatePropagation, set_checked, assign, bind_ended, bind_property, select_option, bubble_event, set_value, effect_root, validate_snippet_args, boundary, wrap_snippet, derived_safe_equal, prop, bind_group, bind_volume, reactive_import, bind_playback_rate, assign_nullish, bind_prop, legacy_pre_effect_reset, remove_input_defaults, aborted, set_selected, validate_store, async_derived, bind_element_size, set_xlink_attribute, trusted, attachment, element, flush, add_svelte_meta, create_ownership_validator, update_pre_store, log_if_contains_state, tag, preventDefault, bind_resize_observer, CLASS, untrack, document, set_default_checked, effect, bind_seekable, snapshot, store_unsub, trace, set_attribute, mark_store_binding, bind_this, set_default_value, remove_textarea_child, NAMESPACE_SVG, slot, attribute_effect, transition, save, set_class, spread_props, HMR, bind_window_scroll, raf, exclude_from_object, init_select, assign_and, for_await_track_reactivity_loss, from_tree, update_prop, set_custom_element_data, bind_buffered, prevent_snippet_stringification, render_effect, bind_files, safe_get, store_get, bind_active_element, tick, await, equals, update_legacy_props, validate_each_keys, proxy, bind_checked, add_locations, async, self, add_legacy_event_listener, bind_current_time, user_pre_effect, validate_dynamic_element_tag, bind_content_editable, update_store, index, invalidate_store, bind_seeking, hydrate_template, clsx, check_target, strict_equals, autofocus, append_styles, noop

// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/constants.js
var constants = __webpack_require__(178);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/index-client.js + 1 modules
var index_client = __webpack_require__(732);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/reactivity/effects.js
var reactivity_effects = __webpack_require__(480);
;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/attachments/index.js
/** @import { Action, ActionReturn } from '../action/public' */ /** @import { Attachment } from './public' */ 



/**
 * Creates an object key that will be recognised as an attachment when the object is spread onto an element,
 * as a programmatic alternative to using `{@attach ...}`. This can be useful for library authors, though
 * is generally not needed when building an app.
 *
 * ```svelte
 * <script>
 * 	import { createAttachmentKey } from 'svelte/attachments';
 *
 * 	const props = {
 * 		class: 'cool',
 * 		onclick: () => alert('clicked'),
 * 		[createAttachmentKey()]: (node) => {
 * 			node.textContent = 'attached!';
 * 		}
 * 	};
 * </script>
 *
 * <button {...props}>click me</button>
 * ```
 * @since 5.29
 */ function createAttachmentKey() {
    return Symbol(ATTACHMENT_KEY);
}
/**
 * Converts an [action](https://svelte.dev/docs/svelte/use) into an [attachment](https://svelte.dev/docs/svelte/@attach) keeping the same behavior.
 * It's useful if you want to start using attachments on components but you have actions provided by a library.
 *
 * Note that the second argument, if provided, must be a function that _returns_ the argument to the
 * action function, not the argument itself.
 *
 * ```svelte
 * <!-- with an action -->
 * <div use:foo={bar}>...</div>
 *
 * <!-- with an attachment -->
 * <div {@attach fromAction(foo, () => bar)}>...</div>
 * ```
 * @template {EventTarget} E
 * @template {unknown} T
 * @overload
 * @param {Action<E, T> | ((element: E, arg: T) => void | ActionReturn<T>)} action The action function
 * @param {() => T} fn A function that returns the argument for the action
 * @returns {Attachment<E>}
 */ /**
 * Converts an [action](https://svelte.dev/docs/svelte/use) into an [attachment](https://svelte.dev/docs/svelte/@attach) keeping the same behavior.
 * It's useful if you want to start using attachments on components but you have actions provided by a library.
 *
 * Note that the second argument, if provided, must be a function that _returns_ the argument to the
 * action function, not the argument itself.
 *
 * ```svelte
 * <!-- with an action -->
 * <div use:foo={bar}>...</div>
 *
 * <!-- with an attachment -->
 * <div {@attach fromAction(foo, () => bar)}>...</div>
 * ```
 * @template {EventTarget} E
 * @overload
 * @param {Action<E, void> | ((element: E) => void | ActionReturn<void>)} action The action function
 * @returns {Attachment<E>}
 */ /**
 * Converts an [action](https://svelte.dev/docs/svelte/use) into an [attachment](https://svelte.dev/docs/svelte/@attach) keeping the same behavior.
 * It's useful if you want to start using attachments on components but you have actions provided by a library.
 *
 * Note that the second argument, if provided, must be a function that _returns_ the argument to the
 * action function, not the argument itself.
 *
 * ```svelte
 * <!-- with an action -->
 * <div use:foo={bar}>...</div>
 *
 * <!-- with an attachment -->
 * <div {@attach fromAction(foo, () => bar)}>...</div>
 * ```
 *
 * @template {EventTarget} E
 * @template {unknown} T
 * @param {Action<E, T> | ((element: E, arg: T) => void | ActionReturn<T>)} action The action function
 * @param {() => T} fn A function that returns the argument for the action
 * @returns {Attachment<E>}
 * @since 5.32
 */ function fromAction(action) {
    let fn = arguments.length > 1 && arguments[1] !== void 0 ? arguments[1] : /** @type {() => T} */ noop;
    return (element)=>{
        const { update, destroy } = untrack(()=>action(element, fn()) ?? {});
        if (update) {
            var ran = false;
            render_effect(()=>{
                const arg = fn();
                if (ran) update(arg);
            });
            ran = true;
        }
        if (destroy) {
            teardown(destroy);
        }
    };
}

// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/context.js
var client_context = __webpack_require__(754);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/utils.js
var utils = __webpack_require__(314);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/runtime.js
var runtime = __webpack_require__(513);
;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dev/assign.js



/**
 *
 * @param {any} a
 * @param {any} b
 * @param {string} property
 * @param {string} location
 */ function compare(a, b, property, location) {
    if (a !== b) {
        w.assignment_value_stale(property, /** @type {string} */ sanitize_location(location));
    }
    return a;
}
/**
 * @param {any} object
 * @param {string} property
 * @param {any} value
 * @param {string} location
 */ function assign_assign(object, property, value, location) {
    return compare(object[property] = value, untrack(()=>object[property]), property, location);
}
/**
 * @param {any} object
 * @param {string} property
 * @param {any} value
 * @param {string} location
 */ function assign_and(object, property, value, location) {
    var _object, _property;
    return compare((_object = object)[_property = property] && (_object[_property] = value), untrack(()=>object[property]), property, location);
}
/**
 * @param {any} object
 * @param {string} property
 * @param {any} value
 * @param {string} location
 */ function assign_or(object, property, value, location) {
    var _object, _property;
    return compare((_object = object)[_property = property] || (_object[_property] = value), untrack(()=>object[property]), property, location);
}
/**
 * @param {any} object
 * @param {string} property
 * @param {any} value
 * @param {string} location
 */ function assign_nullish(object, property, value, location) {
    var _object, _property;
    return compare((_object = object)[_property = property] ?? (_object[_property] = value), untrack(()=>object[property]), property, location);
}

;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dev/css.js
/** @type {Map<String, Set<HTMLStyleElement>>} */ var all_styles = new Map();
/**
 * @param {String} hash
 * @param {HTMLStyleElement} style
 */ function css_register_style(hash, style) {
    var styles = all_styles.get(hash);
    if (!styles) {
        styles = new Set();
        all_styles.set(hash, styles);
    }
    styles.add(style);
}
/**
 * @param {String} hash
 */ function cleanup_styles(hash) {
    var styles = all_styles.get(hash);
    if (!styles) return;
    for (const style of styles){
        style.remove();
    }
    all_styles.delete(hash);
}

// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/constants.js
var client_constants = __webpack_require__(924);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/hydration.js
var hydration = __webpack_require__(452);
;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dev/elements.js
/** @import { SourceLocation } from '#client' */ 



/**
 * @param {any} fn
 * @param {string} filename
 * @param {SourceLocation[]} locations
 * @returns {any}
 */ function add_locations(fn, filename, locations) {
    return function() {
        for(var _len = arguments.length, args = new Array(_len), _key = 0; _key < _len; _key++){
            args[_key] = arguments[_key];
        }
        const dom = fn(...args);
        var node = hydrating ? dom : dom.nodeType === DOCUMENT_FRAGMENT_NODE ? dom.firstChild : dom;
        assign_locations(node, filename, locations);
        return dom;
    };
}
/**
 * @param {Element} element
 * @param {string} filename
 * @param {SourceLocation} location
 */ function assign_location(element, filename, location) {
    // @ts-expect-error
    element.__svelte_meta = {
        parent: dev_stack,
        loc: {
            file: filename,
            line: location[0],
            column: location[1]
        }
    };
    if (location[2]) {
        assign_locations(element.firstChild, filename, location[2]);
    }
}
/**
 * @param {Node | null} node
 * @param {string} filename
 * @param {SourceLocation[]} locations
 */ function assign_locations(node, filename, locations) {
    var i = 0;
    var depth = 0;
    while(node && i < locations.length){
        if (hydrating && node.nodeType === COMMENT_NODE) {
            var comment = /** @type {Comment} */ node;
            if (comment.data === HYDRATION_START || comment.data === HYDRATION_START_ELSE) depth += 1;
            else if (comment.data[0] === HYDRATION_END) depth -= 1;
        }
        if (depth === 0 && node.nodeType === ELEMENT_NODE) {
            assign_location(/** @type {Element} */ node, filename, locations[i++]);
        }
        node = node.nextSibling;
    }
}

// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/reactivity/sources.js
var reactivity_sources = __webpack_require__(264);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/render.js
var render = __webpack_require__(485);
;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dev/hmr.js
/** @import { Source, Effect, TemplateNode } from '#client' */ 






/**
 * @template {(anchor: Comment, props: any) => any} Component
 * @param {Component} original
 * @param {() => Source<Component>} get_source
 */ function hmr(original, get_source) {
    /**
	 * @param {TemplateNode} anchor
	 * @param {any} props
	 */ function wrapper(anchor, props) {
        let instance = {};
        /** @type {Effect} */ let effect;
        let ran = false;
        block(()=>{
            const source = get_source();
            const component = get(source);
            if (effect) {
                // @ts-ignore
                for(var k in instance)delete instance[k];
                destroy_effect(effect);
            }
            effect = branch(()=>{
                // when the component is invalidated, replace it without transitions
                if (ran) set_should_intro(false);
                // preserve getters/setters
                Object.defineProperties(instance, Object.getOwnPropertyDescriptors(// @ts-expect-error
                new.target ? new component(anchor, props) : component(anchor, props)));
                if (ran) set_should_intro(true);
            });
        }, EFFECT_TRANSPARENT);
        ran = true;
        if (hydrating) {
            anchor = hydrate_node;
        }
        return instance;
    }
    // @ts-expect-error
    wrapper[FILENAME] = original[FILENAME];
    // @ts-ignore
    wrapper[HMR] = {
        // When we accept an update, we set the original source to the new component
        original,
        // The `get_source` parameter reads `wrapper[HMR].source`, but in the `accept`
        // function we always replace it with `previous[HMR].source`, which in practice
        // means we only ever update the original
        source: source(original)
    };
    return wrapper;
}

// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/shared/utils.js
var shared_utils = __webpack_require__(986);
;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dev/ownership.js
/** @typedef {{ file: string, line: number, column: number }} Location */ 





/**
 * Sets up a validator that
 * - traverses the path of a prop to find out if it is allowed to be mutated
 * - checks that the binding chain is not interrupted
 * @param {Record<string, any>} props
 */ function create_ownership_validator(props) {
    var _component_context_p;
    const component = component_context === null || component_context === void 0 ? void 0 : component_context.function;
    const parent = component_context === null || component_context === void 0 ? void 0 : (_component_context_p = component_context.p) === null || _component_context_p === void 0 ? void 0 : _component_context_p.function;
    return {
        /**
		 * @param {string} prop
		 * @param {any[]} path
		 * @param {any} result
		 * @param {number} line
		 * @param {number} column
		 */ mutation: (prop, path, result, line, column)=>{
            const name = path[0];
            if (is_bound_or_unset(props, name) || !parent) {
                return result;
            }
            /** @type {any} */ let value = props;
            for(let i = 0; i < path.length - 1; i++){
                value = value[path[i]];
                if (!(value === null || value === void 0 ? void 0 : value[STATE_SYMBOL])) {
                    return result;
                }
            }
            const location = sanitize_location(`${component[FILENAME]}:${line}:${column}`);
            w.ownership_invalid_mutation(name, location, prop, parent[FILENAME]);
            return result;
        },
        /**
		 * @param {any} key
		 * @param {any} child_component
		 * @param {() => any} value
		 */ binding: (key, child_component, value)=>{
            var _value;
            if (!is_bound_or_unset(props, key) && parent && ((_value = value()) === null || _value === void 0 ? void 0 : _value[STATE_SYMBOL])) {
                w.ownership_invalid_binding(component[FILENAME], key, child_component[FILENAME], parent[FILENAME]);
            }
        }
    };
}
/**
 * @param {Record<string, any>} props
 * @param {string} prop_name
 */ function is_bound_or_unset(props, prop_name) {
    var _get_descriptor;
    // Can be the case when someone does `mount(Component, props)` with `let props = $state({...})`
    // or `createClassComponent(Component, props)`
    const is_entry_props = STATE_SYMBOL in props || LEGACY_PROPS in props;
    return !!((_get_descriptor = get_descriptor(props, prop_name)) === null || _get_descriptor === void 0 ? void 0 : _get_descriptor.set) || is_entry_props && prop_name in props || !(prop_name in props);
}

;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dev/legacy.js



/** @param {Function & { [FILENAME]: string }} target */ function check_target(target) {
    if (target) {
        e.component_api_invalid_new(target[FILENAME] ?? 'a component', target.name);
    }
}
function legacy_api() {
    const component = component_context === null || component_context === void 0 ? void 0 : component_context.function;
    /** @param {string} method */ function error(method) {
        e.component_api_changed(method, component[FILENAME]);
    }
    return {
        $destroy: ()=>error('$destroy()'),
        $on: ()=>error('$on(...)'),
        $set: ()=>error('$set(...)')
    };
}

// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dev/tracing.js
var tracing = __webpack_require__(339);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/shared/clone.js
var clone = __webpack_require__(826);
;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dev/inspect.js




/**
 * @param {() => any[]} get_value
 * @param {Function} [inspector]
 */ // eslint-disable-next-line no-console
function inspect(get_value) {
    let inspector = arguments.length > 1 && arguments[1] !== void 0 ? arguments[1] : console.log;
    validate_effect('$inspect');
    let initial = true;
    let error = /** @type {any} */ UNINITIALIZED;
    // Inspect effects runs synchronously so that we can capture useful
    // stack traces. As a consequence, reading the value might result
    // in an error (an `$inspect(object.property)` will run before the
    // `{#if object}...{/if}` that contains it)
    inspect_effect(()=>{
        try {
            var value = get_value();
        } catch (e) {
            error = e;
            return;
        }
        var snap = snapshot(value, true, true);
        untrack(()=>{
            inspector(initial ? 'init' : 'update', ...snap);
        });
        initial = false;
    });
    // If an error occurs, we store it (along with its stack trace).
    // If the render effect subsequently runs, we log the error,
    // but if it doesn't run it's because the `$inspect` was
    // destroyed, meaning we don't need to bother
    render_effect(()=>{
        try {
            // call `get_value` so that this runs alongside the inspect effect
            get_value();
        } catch  {
        // ignore
        }
        if (error !== UNINITIALIZED) {
            // eslint-disable-next-line no-console
            console.error(error);
            error = UNINITIALIZED;
        }
    });
}

// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/reactivity/async.js
var reactivity_async = __webpack_require__(850);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/blocks/boundary.js + 1 modules
var blocks_boundary = __webpack_require__(899);
;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/blocks/async.js
/** @import { TemplateNode, Value } from '#client' */ 



/**
 * @param {TemplateNode} node
 * @param {Array<() => Promise<any>>} expressions
 * @param {(anchor: TemplateNode, ...deriveds: Value[]) => void} fn
 */ function async_async(node, expressions, fn) {
    var boundary = get_boundary();
    boundary.update_pending_count(1);
    var was_hydrating = hydrating;
    if (was_hydrating) {
        hydrate_next();
        var previous_hydrate_node = hydrate_node;
        var end = skip_nodes(false);
        set_hydrate_node(end);
    }
    flatten([], expressions, (values)=>{
        if (was_hydrating) {
            set_hydrating(true);
            set_hydrate_node(previous_hydrate_node);
        }
        try {
            // get values eagerly to avoid creating blocks if they reject
            for (const d of values)get(d);
            fn(node, ...values);
        } finally{
            boundary.update_pending_count(-1);
        }
        if (was_hydrating) {
            set_hydrating(false);
        }
    });
}

// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/task.js
var dom_task = __webpack_require__(593);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/reactivity/batch.js
var reactivity_batch = __webpack_require__(410);
;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/blocks/await.js
/** @import { Effect, Source, TemplateNode } from '#client' */ 









const PENDING = 0;
const THEN = 1;
const CATCH = 2;
/** @typedef {typeof PENDING | typeof THEN | typeof CATCH} AwaitState */ /**
 * @template V
 * @param {TemplateNode} node
 * @param {(() => Promise<V>)} get_input
 * @param {null | ((anchor: Node) => void)} pending_fn
 * @param {null | ((anchor: Node, value: Source<V>) => void)} then_fn
 * @param {null | ((anchor: Node, error: unknown) => void)} catch_fn
 * @returns {void}
 */ function await_block(node, get_input, pending_fn, then_fn, catch_fn) {
    if (hydrating) {
        hydrate_next();
    }
    var anchor = node;
    var runes = is_runes();
    var active_component_context = component_context;
    /** @type {any} */ var component_function = DEV ? component_context === null || component_context === void 0 ? void 0 : component_context.function : null;
    var dev_original_stack = DEV ? dev_stack : null;
    /** @type {V | Promise<V> | typeof UNINITIALIZED} */ var input = UNINITIALIZED;
    /** @type {Effect | null} */ var pending_effect;
    /** @type {Effect | null} */ var then_effect;
    /** @type {Effect | null} */ var catch_effect;
    var input_source = runes ? source(/** @type {V} */ undefined) : mutable_source(/** @type {V} */ undefined, false, false);
    var error_source = runes ? source(undefined) : mutable_source(undefined, false, false);
    var resolved = false;
    /**
	 * @param {AwaitState} state
	 * @param {boolean} restore
	 */ function update(state, restore) {
        resolved = true;
        if (restore) {
            set_active_effect(effect);
            set_active_reaction(effect); // TODO do we need both?
            set_component_context(active_component_context);
            if (DEV) {
                set_dev_current_component_function(component_function);
                set_dev_stack(dev_original_stack);
            }
        }
        try {
            if (state === PENDING && pending_fn) {
                if (pending_effect) resume_effect(pending_effect);
                else pending_effect = branch(()=>pending_fn(anchor));
            }
            if (state === THEN && then_fn) {
                if (then_effect) resume_effect(then_effect);
                else then_effect = branch(()=>then_fn(anchor, input_source));
            }
            if (state === CATCH && catch_fn) {
                if (catch_effect) resume_effect(catch_effect);
                else catch_effect = branch(()=>catch_fn(anchor, error_source));
            }
            if (state !== PENDING && pending_effect) {
                pause_effect(pending_effect, ()=>pending_effect = null);
            }
            if (state !== THEN && then_effect) {
                pause_effect(then_effect, ()=>then_effect = null);
            }
            if (state !== CATCH && catch_effect) {
                pause_effect(catch_effect, ()=>catch_effect = null);
            }
        } finally{
            if (restore) {
                if (DEV) {
                    set_dev_current_component_function(null);
                    set_dev_stack(null);
                }
                set_component_context(null);
                set_active_reaction(null);
                set_active_effect(null);
                // without this, the DOM does not update until two ticks after the promise
                // resolves, which is unexpected behaviour (and somewhat irksome to test)
                flushSync();
            }
        }
    }
    var effect = block(()=>{
        if (input === (input = get_input())) return;
        /** Whether or not there was a hydration mismatch. Needs to be a `let` or else it isn't treeshaken out */ // @ts-ignore coercing `anchor` to a `Comment` causes TypeScript and Prettier to fight
        let mismatch = hydrating && is_promise(input) === (anchor.data === HYDRATION_START_ELSE);
        if (mismatch) {
            // Hydration mismatch: remove everything inside the anchor and start fresh
            anchor = skip_nodes();
            set_hydrate_node(anchor);
            set_hydrating(false);
            mismatch = true;
        }
        if (is_promise(input)) {
            var promise = input;
            resolved = false;
            promise.then((value)=>{
                if (promise !== input) return;
                // we technically could use `set` here since it's on the next microtick
                // but let's use internal_set for consistency and just to be safe
                internal_set(input_source, value);
                update(THEN, true);
            }, (error)=>{
                if (promise !== input) return;
                // we technically could use `set` here since it's on the next microtick
                // but let's use internal_set for consistency and just to be safe
                internal_set(error_source, error);
                update(CATCH, true);
                if (!catch_fn) {
                    // Rethrow the error if no catch block exists
                    throw error_source.v;
                }
            });
            if (hydrating) {
                if (pending_fn) {
                    pending_effect = branch(()=>pending_fn(anchor));
                }
            } else {
                // Wait a microtask before checking if we should show the pending state as
                // the promise might have resolved by the next microtask.
                queue_micro_task(()=>{
                    if (!resolved) update(PENDING, true);
                });
            }
        } else {
            internal_set(input_source, input);
            update(THEN, false);
        }
        if (mismatch) {
            // continue in hydration mode
            set_hydrating(true);
        }
        // Set the input to something else, in order to disable the promise callbacks
        return ()=>input = UNINITIALIZED;
    });
    if (hydrating) {
        anchor = hydrate_node;
    }
}

// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/operations.js
var operations = __webpack_require__(518);
;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/blocks/if.js
/** @import { Effect, TemplateNode } from '#client' */ /** @import { Batch } from '../../reactivity/batch.js'; */ 





// TODO reinstate https://github.com/sveltejs/svelte/pull/15250
/**
 * @param {TemplateNode} node
 * @param {(branch: (fn: (anchor: Node) => void, flag?: boolean) => void) => void} fn
 * @param {boolean} [elseif] True if this is an `{:else if ...}` block rather than an `{#if ...}`, as that affects which transitions are considered 'local'
 * @returns {void}
 */ function if_block(node, fn) {
    let elseif = arguments.length > 2 && arguments[2] !== void 0 ? arguments[2] : false;
    if (hydration/* .hydrating */.fE) {
        (0,hydration/* .hydrate_next */.E$)();
    }
    var anchor = node;
    /** @type {Effect | null} */ var consequent_effect = null;
    /** @type {Effect | null} */ var alternate_effect = null;
    /** @type {typeof UNINITIALIZED | boolean | null} */ var condition = constants/* .UNINITIALIZED */.UP;
    var flags = elseif ? client_constants/* .EFFECT_TRANSPARENT */.lQ : 0;
    var has_branch = false;
    const set_branch = function(/** @type {(anchor: Node) => void} */ fn) {
        let flag = arguments.length > 1 && arguments[1] !== void 0 ? arguments[1] : true;
        has_branch = true;
        update_branch(flag, fn);
    };
    /** @type {DocumentFragment | null} */ var offscreen_fragment = null;
    function commit() {
        if (offscreen_fragment !== null) {
            // remove the anchor
            /** @type {Text} */ offscreen_fragment.lastChild.remove();
            anchor.before(offscreen_fragment);
            offscreen_fragment = null;
        }
        var active = condition ? consequent_effect : alternate_effect;
        var inactive = condition ? alternate_effect : consequent_effect;
        if (active) {
            (0,reactivity_effects/* .resume_effect */.cc)(active);
        }
        if (inactive) {
            (0,reactivity_effects/* .pause_effect */.r4)(inactive, ()=>{
                if (condition) {
                    alternate_effect = null;
                } else {
                    consequent_effect = null;
                }
            });
        }
    }
    const update_branch = (/** @type {boolean | null} */ new_condition, /** @type {null | ((anchor: Node) => void)} */ fn)=>{
        if (condition === (condition = new_condition)) return;
        /** Whether or not there was a hydration mismatch. Needs to be a `let` or else it isn't treeshaken out */ let mismatch = false;
        if (hydration/* .hydrating */.fE) {
            const is_else = (0,hydration/* .read_hydration_instruction */.no)(anchor) === constants/* .HYDRATION_START_ELSE */.qn;
            if (!!condition === is_else) {
                // Hydration mismatch: remove everything inside the anchor and start fresh.
                // This could happen with `{#if browser}...{/if}`, for example
                anchor = (0,hydration/* .skip_nodes */.Ub)();
                (0,hydration/* .set_hydrate_node */.W0)(anchor);
                (0,hydration/* .set_hydrating */.mK)(false);
                mismatch = true;
            }
        }
        var defer = (0,operations/* .should_defer_append */.eL)();
        var target = anchor;
        if (defer) {
            offscreen_fragment = document.createDocumentFragment();
            offscreen_fragment.append(target = (0,operations/* .create_text */.Pb)());
        }
        if (condition) {
            consequent_effect ?? (consequent_effect = fn && (0,reactivity_effects/* .branch */.tk)(()=>fn(target)));
        } else {
            alternate_effect ?? (alternate_effect = fn && (0,reactivity_effects/* .branch */.tk)(()=>fn(target)));
        }
        if (defer) {
            var batch = /** @type {Batch} */ reactivity_batch/* .current_batch */.Dr;
            var active = condition ? consequent_effect : alternate_effect;
            var inactive = condition ? alternate_effect : consequent_effect;
            if (active) batch.skipped_effects.delete(active);
            if (inactive) batch.skipped_effects.add(inactive);
            batch.add_callback(commit);
        } else {
            commit();
        }
        if (mismatch) {
            // continue in hydration mode
            (0,hydration/* .set_hydrating */.mK)(true);
        }
    };
    (0,reactivity_effects/* .block */.om)(()=>{
        has_branch = false;
        fn(set_branch);
        if (!has_branch) {
            update_branch(null, null);
        }
    }, flags);
    if (hydration/* .hydrating */.fE) {
        anchor = hydration/* .hydrate_node */.Xb;
    }
}

;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/blocks/key.js
/** @import { Effect, TemplateNode } from '#client' */ /** @import { Batch } from '../../reactivity/batch.js'; */ 






/**
 * @template V
 * @param {TemplateNode} node
 * @param {() => V} get_key
 * @param {(anchor: Node) => TemplateNode | void} render_fn
 * @returns {void}
 */ function key_key(node, get_key, render_fn) {
    if (hydrating) {
        hydrate_next();
    }
    var anchor = node;
    /** @type {V | typeof UNINITIALIZED} */ var key = UNINITIALIZED;
    /** @type {Effect} */ var effect;
    /** @type {Effect} */ var pending_effect;
    /** @type {DocumentFragment | null} */ var offscreen_fragment = null;
    var changed = is_runes() ? not_equal : safe_not_equal;
    function commit() {
        if (effect) {
            pause_effect(effect);
        }
        if (offscreen_fragment !== null) {
            // remove the anchor
            /** @type {Text} */ offscreen_fragment.lastChild.remove();
            anchor.before(offscreen_fragment);
            offscreen_fragment = null;
        }
        effect = pending_effect;
    }
    block(()=>{
        if (changed(key, key = get_key())) {
            var target = anchor;
            var defer = should_defer_append();
            if (defer) {
                offscreen_fragment = document.createDocumentFragment();
                offscreen_fragment.append(target = create_text());
            }
            pending_effect = branch(()=>render_fn(target));
            if (defer) {
                /** @type {Batch} */ current_batch.add_callback(commit);
            } else {
                commit();
            }
        }
    });
    if (hydrating) {
        anchor = hydrate_node;
    }
}

;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/blocks/css-props.js
/** @import { TemplateNode } from '#client' */ 


/**
 * @param {HTMLDivElement | SVGGElement} element
 * @param {() => Record<string, string>} get_styles
 * @returns {void}
 */ function css_props(element, get_styles) {
    if (hydrating) {
        set_hydrate_node(/** @type {TemplateNode} */ get_first_child(element));
    }
    render_effect(()=>{
        var styles = get_styles();
        for(var key in styles){
            var value = styles[key];
            if (value) {
                element.style.setProperty(key, value);
            } else {
                element.style.removeProperty(key);
            }
        }
    });
}

// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/reactivity/deriveds.js
var deriveds = __webpack_require__(462);
;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/blocks/each.js
/** @import { EachItem, EachState, Effect, MaybeSource, Source, TemplateNode, TransitionManager, Value } from '#client' */ /** @import { Batch } from '../../reactivity/batch.js'; */ 











/**
 * The row of a keyed each block that is currently updating. We track this
 * so that `animate:` directives have something to attach themselves to
 * @type {EachItem | null}
 */ let each_current_each_item = null;
/** @param {EachItem | null} item */ function each_set_current_each_item(item) {
    each_current_each_item = item;
}
/**
 * @param {any} _
 * @param {number} i
 */ function each_index(_, i) {
    return i;
}
/**
 * Pause multiple effects simultaneously, and coordinate their
 * subsequent destruction. Used in each blocks
 * @param {EachState} state
 * @param {EachItem[]} items
 * @param {null | Node} controlled_anchor
 */ function pause_effects(state, items, controlled_anchor) {
    var items_map = state.items;
    /** @type {TransitionManager[]} */ var transitions = [];
    var length = items.length;
    for(var i = 0; i < length; i++){
        pause_children(items[i].e, transitions, true);
    }
    var is_controlled = length > 0 && transitions.length === 0 && controlled_anchor !== null;
    // If we have a controlled anchor, it means that the each block is inside a single
    // DOM element, so we can apply a fast-path for clearing the contents of the element.
    if (is_controlled) {
        var parent_node = /** @type {Element} */ /** @type {Element} */ controlled_anchor.parentNode;
        clear_text_content(parent_node);
        parent_node.append(/** @type {Element} */ controlled_anchor);
        items_map.clear();
        each_link(state, items[0].prev, items[length - 1].next);
    }
    run_out_transitions(transitions, ()=>{
        for(var i = 0; i < length; i++){
            var item = items[i];
            if (!is_controlled) {
                items_map.delete(item.k);
                each_link(state, item.prev, item.next);
            }
            destroy_effect(item.e, !is_controlled);
        }
    });
}
/**
 * @template V
 * @param {Element | Comment} node The next sibling node, or the parent node if this is a 'controlled' block
 * @param {number} flags
 * @param {() => V[]} get_collection
 * @param {(value: V, index: number) => any} get_key
 * @param {(anchor: Node, item: MaybeSource<V>, index: MaybeSource<number>) => void} render_fn
 * @param {null | ((anchor: Node) => void)} fallback_fn
 * @returns {void}
 */ function each(node, flags, get_collection, get_key, render_fn) {
    let fallback_fn = arguments.length > 5 && arguments[5] !== void 0 ? arguments[5] : null;
    var anchor = node;
    /** @type {EachState} */ var state = {
        flags,
        items: new Map(),
        first: null
    };
    var is_controlled = (flags & EACH_IS_CONTROLLED) !== 0;
    if (is_controlled) {
        var parent_node = /** @type {Element} */ node;
        anchor = hydrating ? set_hydrate_node(/** @type {Comment | Text} */ get_first_child(parent_node)) : parent_node.appendChild(create_text());
    }
    if (hydrating) {
        hydrate_next();
    }
    /** @type {Effect | null} */ var fallback = null;
    var was_empty = false;
    /** @type {Map<any, EachItem>} */ var offscreen_items = new Map();
    // TODO: ideally we could use derived for runes mode but because of the ability
    // to use a store which can be mutated, we can't do that here as mutating a store
    // will still result in the collection array being the same from the store
    var each_array = derived_safe_equal(()=>{
        var collection = get_collection();
        return is_array(collection) ? collection : collection == null ? [] : array_from(collection);
    });
    /** @type {V[]} */ var array;
    /** @type {Effect} */ var each_effect;
    function commit() {
        reconcile(each_effect, array, state, offscreen_items, anchor, render_fn, flags, get_key, get_collection);
        if (fallback_fn !== null) {
            if (array.length === 0) {
                if (fallback) {
                    resume_effect(fallback);
                } else {
                    fallback = branch(()=>fallback_fn(anchor));
                }
            } else if (fallback !== null) {
                pause_effect(fallback, ()=>{
                    fallback = null;
                });
            }
        }
    }
    block(()=>{
        // store a reference to the effect so that we can update the start/end nodes in reconciliation
        each_effect ?? (each_effect = /** @type {Effect} */ active_effect);
        array = /** @type {V[]} */ get(each_array);
        var length = array.length;
        if (was_empty && length === 0) {
            // ignore updates if the array is empty,
            // and it already was empty on previous run
            return;
        }
        was_empty = length === 0;
        /** `true` if there was a hydration mismatch. Needs to be a `let` or else it isn't treeshaken out */ let mismatch = false;
        if (hydrating) {
            var is_else = read_hydration_instruction(anchor) === HYDRATION_START_ELSE;
            if (is_else !== (length === 0)) {
                // hydration mismatch — remove the server-rendered DOM and start over
                anchor = skip_nodes();
                set_hydrate_node(anchor);
                set_hydrating(false);
                mismatch = true;
            }
        }
        // this is separate to the previous block because `hydrating` might change
        if (hydrating) {
            /** @type {EachItem | null} */ var prev = null;
            /** @type {EachItem} */ var item;
            for(var i = 0; i < length; i++){
                if (hydrate_node.nodeType === COMMENT_NODE && /** @type {Comment} */ hydrate_node.data === HYDRATION_END) {
                    // The server rendered fewer items than expected,
                    // so break out and continue appending non-hydrated items
                    anchor = /** @type {Comment} */ hydrate_node;
                    mismatch = true;
                    set_hydrating(false);
                    break;
                }
                var value = array[i];
                var key = get_key(value, i);
                item = create_item(hydrate_node, state, prev, null, value, key, i, render_fn, flags, get_collection);
                state.items.set(key, item);
                prev = item;
            }
            // remove excess nodes
            if (length > 0) {
                set_hydrate_node(skip_nodes());
            }
        }
        if (hydrating) {
            if (length === 0 && fallback_fn) {
                fallback = branch(()=>fallback_fn(anchor));
            }
        } else {
            if (should_defer_append()) {
                var keys = new Set();
                var batch = /** @type {Batch} */ current_batch;
                for(i = 0; i < length; i += 1){
                    value = array[i];
                    key = get_key(value, i);
                    var existing = state.items.get(key) ?? offscreen_items.get(key);
                    if (existing) {
                        // update before reconciliation, to trigger any async updates
                        if ((flags & (EACH_ITEM_REACTIVE | EACH_INDEX_REACTIVE)) !== 0) {
                            update_item(existing, value, i, flags);
                        }
                    } else {
                        item = create_item(null, state, null, null, value, key, i, render_fn, flags, get_collection, true);
                        offscreen_items.set(key, item);
                    }
                    keys.add(key);
                }
                for (const [key, item] of state.items){
                    if (!keys.has(key)) {
                        batch.skipped_effects.add(item.e);
                    }
                }
                batch.add_callback(commit);
            } else {
                commit();
            }
        }
        if (mismatch) {
            // continue in hydration mode
            set_hydrating(true);
        }
        // When we mount the each block for the first time, the collection won't be
        // connected to this effect as the effect hasn't finished running yet and its deps
        // won't be assigned. However, it's possible that when reconciling the each block
        // that a mutation occurred and it's made the collection MAYBE_DIRTY, so reading the
        // collection again can provide consistency to the reactive graph again as the deriveds
        // will now be `CLEAN`.
        get(each_array);
    });
    if (hydrating) {
        anchor = hydrate_node;
    }
}
/**
 * Add, remove, or reorder items output by an each block as its input changes
 * @template V
 * @param {Effect} each_effect
 * @param {Array<V>} array
 * @param {EachState} state
 * @param {Map<any, EachItem>} offscreen_items
 * @param {Element | Comment | Text} anchor
 * @param {(anchor: Node, item: MaybeSource<V>, index: number | Source<number>, collection: () => V[]) => void} render_fn
 * @param {number} flags
 * @param {(value: V, index: number) => any} get_key
 * @param {() => V[]} get_collection
 * @returns {void}
 */ function reconcile(each_effect, array, state, offscreen_items, anchor, render_fn, flags, get_key, get_collection) {
    var is_animated = (flags & EACH_IS_ANIMATED) !== 0;
    var should_update = (flags & (EACH_ITEM_REACTIVE | EACH_INDEX_REACTIVE)) !== 0;
    var length = array.length;
    var items = state.items;
    var first = state.first;
    var current = first;
    /** @type {undefined | Set<EachItem>} */ var seen;
    /** @type {EachItem | null} */ var prev = null;
    /** @type {undefined | Set<EachItem>} */ var to_animate;
    /** @type {EachItem[]} */ var matched = [];
    /** @type {EachItem[]} */ var stashed = [];
    /** @type {V} */ var value;
    /** @type {any} */ var key;
    /** @type {EachItem | undefined} */ var item;
    /** @type {number} */ var i;
    if (is_animated) {
        for(i = 0; i < length; i += 1){
            value = array[i];
            key = get_key(value, i);
            item = items.get(key);
            if (item !== undefined) {
                var _item_a;
                (_item_a = item.a) === null || _item_a === void 0 ? void 0 : _item_a.measure();
                (to_animate ?? (to_animate = new Set())).add(item);
            }
        }
    }
    for(i = 0; i < length; i += 1){
        value = array[i];
        key = get_key(value, i);
        item = items.get(key);
        if (item === undefined) {
            var pending = offscreen_items.get(key);
            if (pending !== undefined) {
                offscreen_items.delete(key);
                items.set(key, pending);
                var next = prev ? prev.next : current;
                each_link(state, prev, pending);
                each_link(state, pending, next);
                move(pending, next, anchor);
                prev = pending;
            } else {
                var child_anchor = current ? /** @type {TemplateNode} */ current.e.nodes_start : anchor;
                prev = create_item(child_anchor, state, prev, prev === null ? state.first : prev.next, value, key, i, render_fn, flags, get_collection);
            }
            items.set(key, prev);
            matched = [];
            stashed = [];
            current = prev.next;
            continue;
        }
        if (should_update) {
            update_item(item, value, i, flags);
        }
        if ((item.e.f & INERT) !== 0) {
            resume_effect(item.e);
            if (is_animated) {
                var _item_a1;
                (_item_a1 = item.a) === null || _item_a1 === void 0 ? void 0 : _item_a1.unfix();
                (to_animate ?? (to_animate = new Set())).delete(item);
            }
        }
        if (item !== current) {
            if (seen !== undefined && seen.has(item)) {
                if (matched.length < stashed.length) {
                    // more efficient to move later items to the front
                    var start = stashed[0];
                    var j;
                    prev = start.prev;
                    var a = matched[0];
                    var b = matched[matched.length - 1];
                    for(j = 0; j < matched.length; j += 1){
                        move(matched[j], start, anchor);
                    }
                    for(j = 0; j < stashed.length; j += 1){
                        seen.delete(stashed[j]);
                    }
                    each_link(state, a.prev, b.next);
                    each_link(state, prev, a);
                    each_link(state, b, start);
                    current = start;
                    prev = b;
                    i -= 1;
                    matched = [];
                    stashed = [];
                } else {
                    // more efficient to move earlier items to the back
                    seen.delete(item);
                    move(item, current, anchor);
                    each_link(state, item.prev, item.next);
                    each_link(state, item, prev === null ? state.first : prev.next);
                    each_link(state, prev, item);
                    prev = item;
                }
                continue;
            }
            matched = [];
            stashed = [];
            while(current !== null && current.k !== key){
                // If the each block isn't inert and an item has an effect that is already inert,
                // skip over adding it to our seen Set as the item is already being handled
                if ((current.e.f & INERT) === 0) {
                    (seen ?? (seen = new Set())).add(current);
                }
                stashed.push(current);
                current = current.next;
            }
            if (current === null) {
                continue;
            }
            item = current;
        }
        matched.push(item);
        prev = item;
        current = item.next;
    }
    if (current !== null || seen !== undefined) {
        var to_destroy = seen === undefined ? [] : array_from(seen);
        while(current !== null){
            // If the each block isn't inert, then inert effects are currently outroing and will be removed once the transition is finished
            if ((current.e.f & INERT) === 0) {
                to_destroy.push(current);
            }
            current = current.next;
        }
        var destroy_length = to_destroy.length;
        if (destroy_length > 0) {
            var controlled_anchor = (flags & EACH_IS_CONTROLLED) !== 0 && length === 0 ? anchor : null;
            if (is_animated) {
                for(i = 0; i < destroy_length; i += 1){
                    var _to_destroy_i_a;
                    (_to_destroy_i_a = to_destroy[i].a) === null || _to_destroy_i_a === void 0 ? void 0 : _to_destroy_i_a.measure();
                }
                for(i = 0; i < destroy_length; i += 1){
                    var _to_destroy_i_a1;
                    (_to_destroy_i_a1 = to_destroy[i].a) === null || _to_destroy_i_a1 === void 0 ? void 0 : _to_destroy_i_a1.fix();
                }
            }
            pause_effects(state, to_destroy, controlled_anchor);
        }
    }
    if (is_animated) {
        queue_micro_task(()=>{
            if (to_animate === undefined) return;
            for (item of to_animate){
                var _item_a;
                (_item_a = item.a) === null || _item_a === void 0 ? void 0 : _item_a.apply();
            }
        });
    }
    each_effect.first = state.first && state.first.e;
    each_effect.last = prev && prev.e;
    for (var unused of offscreen_items.values()){
        destroy_effect(unused.e);
    }
    offscreen_items.clear();
}
/**
 * @param {EachItem} item
 * @param {any} value
 * @param {number} index
 * @param {number} type
 * @returns {void}
 */ function update_item(item, value, index, type) {
    if ((type & EACH_ITEM_REACTIVE) !== 0) {
        internal_set(item.v, value);
    }
    if ((type & EACH_INDEX_REACTIVE) !== 0) {
        internal_set(/** @type {Value<number>} */ item.i, index);
    } else {
        item.i = index;
    }
}
/**
 * @template V
 * @param {Node | null} anchor
 * @param {EachState} state
 * @param {EachItem | null} prev
 * @param {EachItem | null} next
 * @param {V} value
 * @param {unknown} key
 * @param {number} index
 * @param {(anchor: Node, item: V | Source<V>, index: number | Value<number>, collection: () => V[]) => void} render_fn
 * @param {number} flags
 * @param {() => V[]} get_collection
 * @param {boolean} [deferred]
 * @returns {EachItem}
 */ function create_item(anchor, state, prev, next, value, key, index, render_fn, flags, get_collection, deferred) {
    var previous_each_item = each_current_each_item;
    var reactive = (flags & EACH_ITEM_REACTIVE) !== 0;
    var mutable = (flags & EACH_ITEM_IMMUTABLE) === 0;
    var v = reactive ? mutable ? mutable_source(value, false, false) : source(value) : value;
    var i = (flags & EACH_INDEX_REACTIVE) === 0 ? index : source(index);
    if (DEV && reactive) {
        // For tracing purposes, we need to link the source signal we create with the
        // collection + index so that tracing works as intended
        /** @type {Value} */ v.trace = ()=>{
            var collection_index = typeof i === 'number' ? index : i.v;
            // eslint-disable-next-line @typescript-eslint/no-unused-expressions
            get_collection()[collection_index];
        };
    }
    /** @type {EachItem} */ var item = {
        i,
        v,
        k: key,
        a: null,
        // @ts-expect-error
        e: null,
        prev,
        next
    };
    each_current_each_item = item;
    try {
        if (anchor === null) {
            var fragment = document.createDocumentFragment();
            fragment.append(anchor = create_text());
        }
        item.e = branch(()=>render_fn(/** @type {Node} */ anchor, v, i, get_collection), hydrating);
        item.e.prev = prev && prev.e;
        item.e.next = next && next.e;
        if (prev === null) {
            if (!deferred) {
                state.first = item;
            }
        } else {
            prev.next = item;
            prev.e.next = item.e;
        }
        if (next !== null) {
            next.prev = item;
            next.e.prev = item.e;
        }
        return item;
    } finally{
        each_current_each_item = previous_each_item;
    }
}
/**
 * @param {EachItem} item
 * @param {EachItem | null} next
 * @param {Text | Element | Comment} anchor
 */ function move(item, next, anchor) {
    var end = item.next ? /** @type {TemplateNode} */ item.next.e.nodes_start : anchor;
    var dest = next ? /** @type {TemplateNode} */ next.e.nodes_start : anchor;
    var node = /** @type {TemplateNode} */ item.e.nodes_start;
    while(node !== null && node !== end){
        var next_node = /** @type {TemplateNode} */ get_next_sibling(node);
        dest.before(node);
        node = next_node;
    }
}
/**
 * @param {EachState} state
 * @param {EachItem | null} prev
 * @param {EachItem | null} next
 */ function each_link(state, prev, next) {
    if (prev === null) {
        state.first = next;
    } else {
        prev.next = next;
        prev.e.next = next && next.e;
    }
    if (next !== null) {
        next.prev = prev;
        next.e.prev = prev && prev.e;
    }
}

// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/reconciler.js
var reconciler = __webpack_require__(642);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/template.js
var template = __webpack_require__(782);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/warnings.js
var warnings = __webpack_require__(32);
// EXTERNAL MODULE: ./node_modules/.pnpm/esm-env@1.2.2/node_modules/esm-env/false.js
var esm_env_false = __webpack_require__(832);
;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/blocks/html.js
/** @import { Effect, TemplateNode } from '#client' */ 











/**
 * @param {Element} element
 * @param {string | null} server_hash
 * @param {string} value
 */ function check_hash(element, server_hash, value) {
    var _element___svelte_meta;
    if (!server_hash || server_hash === (0,utils/* .hash */.tW)(String(value ?? ''))) return;
    let location;
    // @ts-expect-error
    const loc = (_element___svelte_meta = element.__svelte_meta) === null || _element___svelte_meta === void 0 ? void 0 : _element___svelte_meta.loc;
    if (loc) {
        location = `near ${loc.file}:${loc.line}:${loc.column}`;
    } else if (client_context/* .dev_current_component_function */.DE === null || client_context/* .dev_current_component_function */.DE === void 0 ? void 0 : client_context/* .dev_current_component_function */.DE[constants/* .FILENAME */.Uh]) {
        location = `in ${client_context/* .dev_current_component_function */.DE[constants/* .FILENAME */.Uh]}`;
    }
    warnings/* .hydration_html_changed */.Y9((0,utils/* .sanitize_location */.If)(location));
}
/**
 * @param {Element | Text | Comment} node
 * @param {() => string} get_value
 * @param {boolean} [svg]
 * @param {boolean} [mathml]
 * @param {boolean} [skip_warning]
 * @returns {void}
 */ function html_html(node, get_value) {
    let svg = arguments.length > 2 && arguments[2] !== void 0 ? arguments[2] : false, mathml = arguments.length > 3 && arguments[3] !== void 0 ? arguments[3] : false, skip_warning = arguments.length > 4 && arguments[4] !== void 0 ? arguments[4] : false;
    var anchor = node;
    var value = '';
    (0,reactivity_effects/* .template_effect */.vN)(()=>{
        var effect = /** @type {Effect} */ runtime/* .active_effect */.Fg;
        if (value === (value = get_value() ?? '')) {
            if (hydration/* .hydrating */.fE) (0,hydration/* .hydrate_next */.E$)();
            return;
        }
        if (effect.nodes_start !== null) {
            (0,reactivity_effects/* .remove_effect_dom */.mk)(effect.nodes_start, /** @type {TemplateNode} */ effect.nodes_end);
            effect.nodes_start = effect.nodes_end = null;
        }
        if (value === '') return;
        if (hydration/* .hydrating */.fE) {
            // We're deliberately not trying to repair mismatches between server and client,
            // as it's costly and error-prone (and it's an edge case to have a mismatch anyway)
            var hash = /** @type {Comment} */ hydration/* .hydrate_node.data */.Xb.data;
            var next = (0,hydration/* .hydrate_next */.E$)();
            var last = next;
            while(next !== null && (next.nodeType !== client_constants/* .COMMENT_NODE */.dz || /** @type {Comment} */ next.data !== '')){
                last = next;
                next = /** @type {TemplateNode} */ (0,operations/* .get_next_sibling */.M$)(next);
            }
            if (next === null) {
                warnings/* .hydration_mismatch */.eZ();
                throw constants/* .HYDRATION_ERROR */.kD;
            }
            if (esm_env_false/* ["default"] */.A && !skip_warning) {
                check_hash(/** @type {Element} */ next.parentNode, hash, value);
            }
            (0,template/* .assign_nodes */.mX)(hydration/* .hydrate_node */.Xb, last);
            anchor = (0,hydration/* .set_hydrate_node */.W0)(next);
            return;
        }
        var html = value + '';
        if (svg) html = `<svg>${html}</svg>`;
        else if (mathml) html = `<math>${html}</math>`;
        // Don't use create_fragment_with_script_from_html here because that would mean script tags are executed.
        // @html is basically `.innerHTML = ...` and that doesn't execute scripts either due to security reasons.
        /** @type {DocumentFragment | Element} */ var node = (0,reconciler/* .create_fragment_from_html */.L)(html);
        if (svg || mathml) {
            node = /** @type {Element} */ (0,operations/* .get_first_child */.Zj)(node);
        }
        (0,template/* .assign_nodes */.mX)(/** @type {TemplateNode} */ (0,operations/* .get_first_child */.Zj)(node), /** @type {TemplateNode} */ node.lastChild);
        if (svg || mathml) {
            while((0,operations/* .get_first_child */.Zj)(node)){
                anchor.before(/** @type {Node} */ (0,operations/* .get_first_child */.Zj)(node));
            }
        } else {
            anchor.before(node);
        }
    });
}

;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/blocks/slot.js

/**
 * @param {Comment} anchor
 * @param {Record<string, any>} $$props
 * @param {string} name
 * @param {Record<string, unknown>} slot_props
 * @param {null | ((anchor: Comment) => void)} fallback_fn
 */ function slot_slot(anchor, $$props, name, slot_props, fallback_fn) {
    var _$$props_$$slots;
    if (hydrating) {
        hydrate_next();
    }
    var slot_fn = (_$$props_$$slots = $$props.$$slots) === null || _$$props_$$slots === void 0 ? void 0 : _$$props_$$slots[name];
    // Interop: Can use snippets to fill slots
    var is_interop = false;
    if (slot_fn === true) {
        slot_fn = $$props[name === 'default' ? 'children' : name];
        is_interop = true;
    }
    if (slot_fn === undefined) {
        if (fallback_fn !== null) {
            fallback_fn(anchor);
        }
    } else {
        slot_fn(anchor, is_interop ? ()=>slot_props : slot_props);
    }
}
/**
 * @param {Record<string, any>} props
 * @returns {Record<string, boolean>}
 */ function sanitize_slots(props) {
    /** @type {Record<string, boolean>} */ const sanitized = {};
    if (props.children) sanitized.default = true;
    for(const key in props.$$slots){
        sanitized[key] = true;
    }
    return sanitized;
}

// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/blocks/snippet.js
var snippet = __webpack_require__(768);
;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/blocks/svelte-component.js
/** @import { TemplateNode, Dom, Effect } from '#client' */ /** @import { Batch } from '../../reactivity/batch.js'; */ 




/**
 * @template P
 * @template {(props: P) => void} C
 * @param {TemplateNode} node
 * @param {() => C} get_component
 * @param {(anchor: TemplateNode, component: C) => Dom | void} render_fn
 * @returns {void}
 */ function svelte_component_component(node, get_component, render_fn) {
    if (hydrating) {
        hydrate_next();
    }
    var anchor = node;
    /** @type {C} */ var component;
    /** @type {Effect | null} */ var effect;
    /** @type {DocumentFragment | null} */ var offscreen_fragment = null;
    /** @type {Effect | null} */ var pending_effect = null;
    function commit() {
        if (effect) {
            pause_effect(effect);
            effect = null;
        }
        if (offscreen_fragment) {
            // remove the anchor
            /** @type {Text} */ offscreen_fragment.lastChild.remove();
            anchor.before(offscreen_fragment);
            offscreen_fragment = null;
        }
        effect = pending_effect;
        pending_effect = null;
    }
    block(()=>{
        if (component === (component = get_component())) return;
        var defer = should_defer_append();
        if (component) {
            var target = anchor;
            if (defer) {
                offscreen_fragment = document.createDocumentFragment();
                offscreen_fragment.append(target = create_text());
                if (effect) {
                    /** @type {Batch} */ current_batch.skipped_effects.add(effect);
                }
            }
            pending_effect = branch(()=>render_fn(target, component));
        }
        if (defer) {
            /** @type {Batch} */ current_batch.add_callback(commit);
        } else {
            commit();
        }
    }, EFFECT_TRANSPARENT);
    if (hydrating) {
        anchor = hydrate_node;
    }
}

;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/blocks/svelte-element.js
/** @import { Effect, TemplateNode } from '#client' */ 











/**
 * @param {Comment | Element} node
 * @param {() => string} get_tag
 * @param {boolean} is_svg
 * @param {undefined | ((element: Element, anchor: Node | null) => void)} render_fn,
 * @param {undefined | (() => string)} get_namespace
 * @param {undefined | [number, number]} location
 * @returns {void}
 */ function svelte_element_element(node, get_tag, is_svg, render_fn, get_namespace, location) {
    let was_hydrating = hydrating;
    if (hydrating) {
        hydrate_next();
    }
    var filename = DEV && location && (component_context === null || component_context === void 0 ? void 0 : component_context.function[FILENAME]);
    /** @type {string | null} */ var tag;
    /** @type {string | null} */ var current_tag;
    /** @type {null | Element} */ var element = null;
    if (hydrating && hydrate_node.nodeType === ELEMENT_NODE) {
        element = /** @type {Element} */ hydrate_node;
        hydrate_next();
    }
    var anchor = /** @type {TemplateNode} */ hydrating ? hydrate_node : node;
    /** @type {Effect | null} */ var effect;
    /**
	 * The keyed `{#each ...}` item block, if any, that this element is inside.
	 * We track this so we can set it when changing the element, allowing any
	 * `animate:` directive to bind itself to the correct block
	 */ var each_item_block = current_each_item;
    block(()=>{
        const next_tag = get_tag() || null;
        var ns = get_namespace ? get_namespace() : is_svg || next_tag === 'svg' ? NAMESPACE_SVG : null;
        // Assumption: Noone changes the namespace but not the tag (what would that even mean?)
        if (next_tag === tag) return;
        // See explanation of `each_item_block` above
        var previous_each_item = current_each_item;
        set_current_each_item(each_item_block);
        if (effect) {
            if (next_tag === null) {
                // start outro
                pause_effect(effect, ()=>{
                    effect = null;
                    current_tag = null;
                });
            } else if (next_tag === current_tag) {
                // same tag as is currently rendered — abort outro
                resume_effect(effect);
            } else {
                // tag is changing — destroy immediately, render contents without intro transitions
                destroy_effect(effect);
                set_should_intro(false);
            }
        }
        if (next_tag && next_tag !== current_tag) {
            effect = branch(()=>{
                element = hydrating ? /** @type {Element} */ element : ns ? document.createElementNS(ns, next_tag) : document.createElement(next_tag);
                if (DEV && location) {
                    // @ts-expect-error
                    element.__svelte_meta = {
                        parent: dev_stack,
                        loc: {
                            file: filename,
                            line: location[0],
                            column: location[1]
                        }
                    };
                }
                assign_nodes(element, element);
                if (render_fn) {
                    if (hydrating && is_raw_text_element(next_tag)) {
                        // prevent hydration glitches
                        element.append(document.createComment(''));
                    }
                    // If hydrating, use the existing ssr comment as the anchor so that the
                    // inner open and close methods can pick up the existing nodes correctly
                    var child_anchor = /** @type {TemplateNode} */ hydrating ? get_first_child(element) : element.appendChild(create_text());
                    if (hydrating) {
                        if (child_anchor === null) {
                            set_hydrating(false);
                        } else {
                            set_hydrate_node(child_anchor);
                        }
                    }
                    // `child_anchor` is undefined if this is a void element, but we still
                    // need to call `render_fn` in order to run actions etc. If the element
                    // contains children, it's a user error (which is warned on elsewhere)
                    // and the DOM will be silently discarded
                    render_fn(element, child_anchor);
                }
                // we do this after calling `render_fn` so that child effects don't override `nodes.end`
                /** @type {Effect} */ active_effect.nodes_end = element;
                anchor.before(element);
            });
        }
        tag = next_tag;
        if (tag) current_tag = tag;
        set_should_intro(true);
        set_current_each_item(previous_each_item);
    }, EFFECT_TRANSPARENT);
    if (was_hydrating) {
        set_hydrating(true);
        set_hydrate_node(anchor);
    }
}

// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/blocks/svelte-head.js
var svelte_head = __webpack_require__(777);
;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/css.js



/**
 * @param {Node} anchor
 * @param {{ hash: string, code: string }} css
 */ function append_styles(anchor, css) {
    // Use `queue_micro_task` to ensure `anchor` is in the DOM, otherwise getRootNode() will yield wrong results
    effect(()=>{
        var root = anchor.getRootNode();
        var target = /** @type {ShadowRoot} */ root.host ? /** @type {ShadowRoot} */ root : /** @type {Document} */ root.head ?? /** @type {Document} */ root.ownerDocument.head;
        // Always querying the DOM is roughly the same perf as additionally checking for presence in a map first assuming
        // that you'll get cache hits half of the time, so we just always query the dom for simplicity and code savings.
        if (!target.querySelector('#' + css.hash)) {
            const style = document.createElement('style');
            style.id = css.hash;
            style.textContent = css.code;
            target.appendChild(style);
            if (DEV) {
                register_style(css.hash, style);
            }
        }
    });
}

// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/reactivity/equality.js
var equality = __webpack_require__(576);
;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/elements/actions.js
/** @import { ActionPayload } from '#client' */ 


/**
 * @template P
 * @param {Element} dom
 * @param {(dom: Element, value?: P) => ActionPayload<P>} action
 * @param {() => P} [get_value]
 * @returns {void}
 */ function actions_action(dom, action, get_value) {
    (0,reactivity_effects/* .effect */.QZ)(()=>{
        var payload = (0,runtime/* .untrack */.vz)(()=>action(dom, get_value === null || get_value === void 0 ? void 0 : get_value()) || {});
        if (get_value && (payload === null || payload === void 0 ? void 0 : payload.update)) {
            var inited = false;
            /** @type {P} */ var prev = /** @type {any} */ {}; // initialize with something so it's never equal on first run
            (0,reactivity_effects/* .render_effect */.VB)(()=>{
                var value = get_value();
                // Action's update method is coarse-grained, i.e. when anything in the passed value changes, update.
                // This works in legacy mode because of mutable_source being updated as a whole, but when using $state
                // together with actions and mutation, it wouldn't notice the change without a deep read.
                (0,runtime/* .deep_read_state */.iT)(value);
                if (inited && (0,equality/* .safe_not_equal */.jX)(prev, value)) {
                    prev = value;
                    /** @type {Function} */ payload.update(value);
                }
            });
            inited = true;
        }
        if (payload === null || payload === void 0 ? void 0 : payload.destroy) {
            return ()=>/** @type {Function} */ payload.destroy();
        }
    });
}

;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/elements/attachments.js
/** @import { Effect } from '#client' */ 
// TODO in 6.0 or 7.0, when we remove legacy mode, we can simplify this by
// getting rid of the block/branch stuff and just letting the effect rip.
// see https://github.com/sveltejs/svelte/pull/15962
/**
 * @param {Element} node
 * @param {() => (node: Element) => void} get_fn
 */ function attachments_attach(node, get_fn) {
    /** @type {false | undefined | ((node: Element) => void)} */ var fn = undefined;
    /** @type {Effect | null} */ var e;
    block(()=>{
        if (fn !== (fn = get_fn())) {
            if (e) {
                destroy_effect(e);
                e = null;
            }
            if (fn) {
                e = branch(()=>{
                    effect(()=>/** @type {(node: Element) => void} */ fn(node));
                });
            }
        }
    });
}

// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/elements/events.js
var elements_events = __webpack_require__(417);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/elements/misc.js
var misc = __webpack_require__(108);
;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/shared/attributes.js


/**
 * `<div translate={false}>` should be rendered as `<div translate="no">` and _not_
 * `<div translate="false">`, which is equivalent to `<div translate="yes">`. There
 * may be other odd cases that need to be added to this list in future
 * @type {Record<string, Map<any, string>>}
 */ const replacements = {
    translate: new Map([
        [
            true,
            'yes'
        ],
        [
            false,
            'no'
        ]
    ])
};
/**
 * @template V
 * @param {string} name
 * @param {V} value
 * @param {boolean} [is_boolean]
 * @returns {string}
 */ function attributes_attr(name, value) {
    let is_boolean = arguments.length > 2 && arguments[2] !== void 0 ? arguments[2] : false;
    // attribute hidden for values other than "until-found" behaves like a boolean attribute
    if (name === 'hidden' && value !== 'until-found') {
        is_boolean = true;
    }
    if (value == null || !value && is_boolean) return '';
    const normalized = name in replacements && replacements[name].get(value) || value;
    const assignment = is_boolean ? '' : `="${escape_html(normalized, true)}"`;
    return ` ${name}${assignment}`;
}
/**
 * Small wrapper around clsx to preserve Svelte's (weird) handling of falsy values.
 * TODO Svelte 6 revisit this, and likely turn all falsy values into the empty string (what clsx also does)
 * @param  {any} value
 */ function attributes_clsx(value) {
    if (typeof value === 'object') {
        return _clsx(value);
    } else {
        return value ?? '';
    }
}
const whitespace = [
    ...' \t\n\r\f\u00a0\u000b\ufeff'
];
/**
 * @param {any} value
 * @param {string | null} [hash]
 * @param {Record<string, boolean>} [directives]
 * @returns {string | null}
 */ function attributes_to_class(value, hash, directives) {
    var classname = value == null ? '' : '' + value;
    if (hash) {
        classname = classname ? classname + ' ' + hash : hash;
    }
    if (directives) {
        for(var key in directives){
            if (directives[key]) {
                classname = classname ? classname + ' ' + key : key;
            } else if (classname.length) {
                var len = key.length;
                var a = 0;
                while((a = classname.indexOf(key, a)) >= 0){
                    var b = a + len;
                    if ((a === 0 || whitespace.includes(classname[a - 1])) && (b === classname.length || whitespace.includes(classname[b]))) {
                        classname = (a === 0 ? '' : classname.substring(0, a)) + classname.substring(b + 1);
                    } else {
                        a = b;
                    }
                }
            }
        }
    }
    return classname === '' ? null : classname;
}
/**
 *
 * @param {Record<string,any>} styles
 * @param {boolean} important
 */ function attributes_append_styles(styles) {
    let important = arguments.length > 1 && arguments[1] !== void 0 ? arguments[1] : false;
    var separator = important ? ' !important;' : ';';
    var css = '';
    for(var key in styles){
        var value = styles[key];
        if (value != null && value !== '') {
            css += ' ' + key + ': ' + value + separator;
        }
    }
    return css;
}
/**
 * @param {string} name
 * @returns {string}
 */ function to_css_name(name) {
    if (name[0] !== '-' || name[1] !== '-') {
        return name.toLowerCase();
    }
    return name;
}
/**
 * @param {any} value
 * @param {Record<string, any> | [Record<string, any>, Record<string, any>]} [styles]
 * @returns {string | null}
 */ function attributes_to_style(value, styles) {
    if (styles) {
        var new_style = '';
        /** @type {Record<string,any> | undefined} */ var normal_styles;
        /** @type {Record<string,any> | undefined} */ var important_styles;
        if (Array.isArray(styles)) {
            normal_styles = styles[0];
            important_styles = styles[1];
        } else {
            normal_styles = styles;
        }
        if (value) {
            value = String(value).replaceAll(/\s*\/\*.*?\*\/\s*/g, '').trim();
            /** @type {boolean | '"' | "'"} */ var in_str = false;
            var in_apo = 0;
            var in_comment = false;
            var reserved_names = [];
            if (normal_styles) {
                reserved_names.push(...Object.keys(normal_styles).map(to_css_name));
            }
            if (important_styles) {
                reserved_names.push(...Object.keys(important_styles).map(to_css_name));
            }
            var start_index = 0;
            var name_index = -1;
            const len = value.length;
            for(var i = 0; i < len; i++){
                var c = value[i];
                if (in_comment) {
                    if (c === '/' && value[i - 1] === '*') {
                        in_comment = false;
                    }
                } else if (in_str) {
                    if (in_str === c) {
                        in_str = false;
                    }
                } else if (c === '/' && value[i + 1] === '*') {
                    in_comment = true;
                } else if (c === '"' || c === "'") {
                    in_str = c;
                } else if (c === '(') {
                    in_apo++;
                } else if (c === ')') {
                    in_apo--;
                }
                if (!in_comment && in_str === false && in_apo === 0) {
                    if (c === ':' && name_index === -1) {
                        name_index = i;
                    } else if (c === ';' || i === len - 1) {
                        if (name_index !== -1) {
                            var name = to_css_name(value.substring(start_index, name_index).trim());
                            if (!reserved_names.includes(name)) {
                                if (c !== ';') {
                                    i++;
                                }
                                var property = value.substring(start_index, i).trim();
                                new_style += ' ' + property + ';';
                            }
                        }
                        start_index = i + 1;
                        name_index = -1;
                    }
                }
            }
        }
        if (normal_styles) {
            new_style += attributes_append_styles(normal_styles);
        }
        if (important_styles) {
            new_style += attributes_append_styles(important_styles, true);
        }
        new_style = new_style.trim();
        return new_style === '' ? null : new_style;
    }
    return value == null ? null : String(value);
}

;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/elements/class.js


/**
 * @param {Element} dom
 * @param {boolean | number} is_html
 * @param {string | null} value
 * @param {string} [hash]
 * @param {Record<string, any>} [prev_classes]
 * @param {Record<string, any>} [next_classes]
 * @returns {Record<string, boolean> | undefined}
 */ function class_set_class(dom, is_html, value, hash, prev_classes, next_classes) {
    // @ts-expect-error need to add __className to patched prototype
    var prev = dom.__className;
    if (hydrating || prev !== value || prev === undefined // for edge case of `class={undefined}`
    ) {
        var next_class_name = to_class(value, hash, next_classes);
        if (!hydrating || next_class_name !== dom.getAttribute('class')) {
            // Removing the attribute when the value is only an empty string causes
            // performance issues vs simply making the className an empty string. So
            // we should only remove the class if the value is nullish
            // and there no hash/directives :
            if (next_class_name == null) {
                dom.removeAttribute('class');
            } else if (is_html) {
                dom.className = next_class_name;
            } else {
                dom.setAttribute('class', next_class_name);
            }
        }
        // @ts-expect-error need to add __className to patched prototype
        dom.__className = value;
    } else if (next_classes && prev_classes !== next_classes) {
        for(var key in next_classes){
            var is_present = !!next_classes[key];
            if (prev_classes == null || is_present !== !!prev_classes[key]) {
                dom.classList.toggle(key, is_present);
            }
        }
    }
    return next_classes;
}

;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/elements/style.js


/**
 * @param {Element & ElementCSSInlineStyle} dom
 * @param {Record<string, any>} prev
 * @param {Record<string, any>} next
 * @param {string} [priority]
 */ function update_styles(dom) {
    let prev = arguments.length > 1 && arguments[1] !== void 0 ? arguments[1] : {}, next = arguments.length > 2 ? arguments[2] : void 0, priority = arguments.length > 3 ? arguments[3] : void 0;
    for(var key in next){
        var value = next[key];
        if (prev[key] !== value) {
            if (next[key] == null) {
                dom.style.removeProperty(key);
            } else {
                dom.style.setProperty(key, value, priority);
            }
        }
    }
}
/**
 * @param {Element & ElementCSSInlineStyle} dom
 * @param {string | null} value
 * @param {Record<string, any> | [Record<string, any>, Record<string, any>]} [prev_styles]
 * @param {Record<string, any> | [Record<string, any>, Record<string, any>]} [next_styles]
 */ function style_set_style(dom, value, prev_styles, next_styles) {
    // @ts-expect-error
    var prev = dom.__style;
    if (hydrating || prev !== value) {
        var next_style_attr = to_style(value, next_styles);
        if (!hydrating || next_style_attr !== dom.getAttribute('style')) {
            if (next_style_attr == null) {
                dom.removeAttribute('style');
            } else {
                dom.style.cssText = next_style_attr;
            }
        }
        // @ts-expect-error
        dom.__style = value;
    } else if (next_styles) {
        if (Array.isArray(next_styles)) {
            update_styles(dom, prev_styles === null || prev_styles === void 0 ? void 0 : prev_styles[0], next_styles[0]);
            update_styles(dom, prev_styles === null || prev_styles === void 0 ? void 0 : prev_styles[1], next_styles[1], 'important');
        } else {
            update_styles(dom, prev_styles, next_styles);
        }
    }
    return next_styles;
}

// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/elements/bindings/shared.js
var shared = __webpack_require__(408);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/proxy.js
var client_proxy = __webpack_require__(445);
;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/elements/bindings/select.js





/**
 * Selects the correct option(s) (depending on whether this is a multiple select)
 * @template V
 * @param {HTMLSelectElement} select
 * @param {V} value
 * @param {boolean} mounting
 */ function select_select_option(select, value) {
    let mounting = arguments.length > 2 && arguments[2] !== void 0 ? arguments[2] : false;
    if (select.multiple) {
        // If value is null or undefined, keep the selection as is
        if (value == undefined) {
            return;
        }
        // If not an array, warn and keep the selection as is
        if (!is_array(value)) {
            return w.select_multiple_invalid_value();
        }
        // Otherwise, update the selection
        for (var option of select.options){
            option.selected = value.includes(get_option_value(option));
        }
        return;
    }
    for (option of select.options){
        var option_value = get_option_value(option);
        if (is(option_value, value)) {
            option.selected = true;
            return;
        }
    }
    if (!mounting || value !== undefined) {
        select.selectedIndex = -1; // no option should be selected
    }
}
/**
 * Selects the correct option(s) if `value` is given,
 * and then sets up a mutation observer to sync the
 * current selection to the dom when it changes. Such
 * changes could for example occur when options are
 * inside an `#each` block.
 * @param {HTMLSelectElement} select
 */ function select_init_select(select) {
    var observer = new MutationObserver(()=>{
        // @ts-ignore
        select_select_option(select, select.__value);
    // Deliberately don't update the potential binding value,
    // the model should be preserved unless explicitly changed
    });
    observer.observe(select, {
        // Listen to option element changes
        childList: true,
        subtree: true,
        // Listen to option element value attribute changes
        // (doesn't get notified of select value changes,
        // because that property is not reflected as an attribute)
        attributes: true,
        attributeFilter: [
            'value'
        ]
    });
    teardown(()=>{
        observer.disconnect();
    });
}
/**
 * @param {HTMLSelectElement} select
 * @param {() => unknown} get
 * @param {(value: unknown) => void} set
 * @returns {void}
 */ function bind_select_value(select, get) {
    let set = arguments.length > 2 && arguments[2] !== void 0 ? arguments[2] : get;
    var mounting = true;
    listen_to_event_and_reset_event(select, 'change', (is_reset)=>{
        var query = is_reset ? '[selected]' : ':checked';
        /** @type {unknown} */ var value;
        if (select.multiple) {
            value = [].map.call(select.querySelectorAll(query), get_option_value);
        } else {
            /** @type {HTMLOptionElement | null} */ var selected_option = select.querySelector(query) ?? // will fall back to first non-disabled option if no option is selected
            select.querySelector('option:not([disabled])');
            value = selected_option && get_option_value(selected_option);
        }
        set(value);
    });
    // Needs to be an effect, not a render_effect, so that in case of each loops the logic runs after the each block has updated
    effect(()=>{
        var value = get();
        select_select_option(select, value, mounting);
        // Mounting and value undefined -> take selection from dom
        if (mounting && value === undefined) {
            /** @type {HTMLOptionElement | null} */ var selected_option = select.querySelector(':checked');
            if (selected_option !== null) {
                value = get_option_value(selected_option);
                set(value);
            }
        }
        // @ts-ignore
        select.__value = value;
        mounting = false;
    });
    select_init_select(select);
}
/** @param {HTMLOptionElement} option */ function get_option_value(option) {
    // __value only exists if the <option> has a value attribute
    if ('__value' in option) {
        return option.__value;
    } else {
        return option.value;
    }
}

;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/elements/attributes.js
/** @import { Effect } from '#client' */ 

















const CLASS = Symbol('class');
const STYLE = Symbol('style');
const IS_CUSTOM_ELEMENT = Symbol('is custom element');
const IS_HTML = Symbol('is html');
/**
 * The value/checked attribute in the template actually corresponds to the defaultValue property, so we need
 * to remove it upon hydration to avoid a bug when someone resets the form value.
 * @param {HTMLInputElement} input
 * @returns {void}
 */ function remove_input_defaults(input) {
    if (!hydrating) return;
    var already_removed = false;
    // We try and remove the default attributes later, rather than sync during hydration.
    // Doing it sync during hydration has a negative impact on performance, but deferring the
    // work in an idle task alleviates this greatly. If a form reset event comes in before
    // the idle callback, then we ensure the input defaults are cleared just before.
    var remove_defaults = ()=>{
        if (already_removed) return;
        already_removed = true;
        // Remove the attributes but preserve the values
        if (input.hasAttribute('value')) {
            var value = input.value;
            set_attribute(input, 'value', null);
            input.value = value;
        }
        if (input.hasAttribute('checked')) {
            var checked = input.checked;
            set_attribute(input, 'checked', null);
            input.checked = checked;
        }
    };
    // @ts-expect-error
    input.__on_r = remove_defaults;
    queue_idle_task(remove_defaults);
    add_form_reset_listener();
}
/**
 * @param {Element} element
 * @param {any} value
 */ function set_value(element, value) {
    var attributes = get_attributes(element);
    if (attributes.value === (attributes.value = // treat null and undefined the same for the initial value
    value ?? undefined) || // @ts-expect-error
    // `progress` elements always need their value set when it's `0`
    element.value === value && (value !== 0 || element.nodeName !== 'PROGRESS')) {
        return;
    }
    // @ts-expect-error
    element.value = value ?? '';
}
/**
 * @param {Element} element
 * @param {boolean} checked
 */ function set_checked(element, checked) {
    var attributes = get_attributes(element);
    if (attributes.checked === (attributes.checked = // treat null and undefined the same for the initial value
    checked ?? undefined)) {
        return;
    }
    // @ts-expect-error
    element.checked = checked;
}
/**
 * Sets the `selected` attribute on an `option` element.
 * Not set through the property because that doesn't reflect to the DOM,
 * which means it wouldn't be taken into account when a form is reset.
 * @param {HTMLOptionElement} element
 * @param {boolean} selected
 */ function set_selected(element, selected) {
    if (selected) {
        // The selected option could've changed via user selection, and
        // setting the value without this check would set it back.
        if (!element.hasAttribute('selected')) {
            element.setAttribute('selected', '');
        }
    } else {
        element.removeAttribute('selected');
    }
}
/**
 * Applies the default checked property without influencing the current checked property.
 * @param {HTMLInputElement} element
 * @param {boolean} checked
 */ function set_default_checked(element, checked) {
    const existing_value = element.checked;
    element.defaultChecked = checked;
    element.checked = existing_value;
}
/**
 * Applies the default value property without influencing the current value property.
 * @param {HTMLInputElement | HTMLTextAreaElement} element
 * @param {string} value
 */ function set_default_value(element, value) {
    const existing_value = element.value;
    element.defaultValue = value;
    element.value = existing_value;
}
/**
 * @param {Element} element
 * @param {string} attribute
 * @param {string | null} value
 * @param {boolean} [skip_warning]
 */ function set_attribute(element, attribute, value, skip_warning) {
    var attributes = get_attributes(element);
    if (hydrating) {
        attributes[attribute] = element.getAttribute(attribute);
        if (attribute === 'src' || attribute === 'srcset' || attribute === 'href' && element.nodeName === 'LINK') {
            if (!skip_warning) {
                check_src_in_dev_hydration(element, attribute, value ?? '');
            }
            // If we reset these attributes, they would result in another network request, which we want to avoid.
            // We assume they are the same between client and server as checking if they are equal is expensive
            // (we can't just compare the strings as they can be different between client and server but result in the
            // same url, so we would need to create hidden anchor elements to compare them)
            return;
        }
    }
    if (attributes[attribute] === (attributes[attribute] = value)) return;
    if (attribute === 'loading') {
        // @ts-expect-error
        element[LOADING_ATTR_SYMBOL] = value;
    }
    if (value == null) {
        element.removeAttribute(attribute);
    } else if (typeof value !== 'string' && get_setters(element).includes(attribute)) {
        // @ts-ignore
        element[attribute] = value;
    } else {
        element.setAttribute(attribute, value);
    }
}
/**
 * @param {Element} dom
 * @param {string} attribute
 * @param {string} value
 */ function set_xlink_attribute(dom, attribute, value) {
    dom.setAttributeNS('http://www.w3.org/1999/xlink', attribute, value);
}
/**
 * @param {HTMLElement} node
 * @param {string} prop
 * @param {any} value
 */ function set_custom_element_data(node, prop, value) {
    // We need to ensure that setting custom element props, which can
    // invoke lifecycle methods on other custom elements, does not also
    // associate those lifecycle methods with the current active reaction
    // or effect
    var previous_reaction = active_reaction;
    var previous_effect = active_effect;
    // If we're hydrating but the custom element is from Svelte, and it already scaffolded,
    // then it might run block logic in hydration mode, which we have to prevent.
    let was_hydrating = hydrating;
    if (hydrating) {
        set_hydrating(false);
    }
    set_active_reaction(null);
    set_active_effect(null);
    try {
        if (// `style` should use `set_attribute` rather than the setter
        prop !== 'style' && // Don't compute setters for custom elements while they aren't registered yet,
        // because during their upgrade/instantiation they might add more setters.
        // Instead, fall back to a simple "an object, then set as property" heuristic.
        (setters_cache.has(node.getAttribute('is') || node.nodeName) || // customElements may not be available in browser extension contexts
        !customElements || customElements.get(node.getAttribute('is') || node.tagName.toLowerCase()) ? get_setters(node).includes(prop) : value && typeof value === 'object')) {
            // @ts-expect-error
            node[prop] = value;
        } else {
            // We did getters etc checks already, stringify before passing to set_attribute
            // to ensure it doesn't invoke the same logic again, and potentially populating
            // the setters cache too early.
            set_attribute(node, prop, value == null ? value : String(value));
        }
    } finally{
        set_active_reaction(previous_reaction);
        set_active_effect(previous_effect);
        if (was_hydrating) {
            set_hydrating(true);
        }
    }
}
/**
 * Spreads attributes onto a DOM element, taking into account the currently set attributes
 * @param {Element & ElementCSSInlineStyle} element
 * @param {Record<string | symbol, any> | undefined} prev
 * @param {Record<string | symbol, any>} next New attributes - this function mutates this object
 * @param {string} [css_hash]
 * @param {boolean} [should_remove_defaults]
 * @param {boolean} [skip_warning]
 * @returns {Record<string, any>}
 */ function set_attributes(element, prev, next, css_hash) {
    let should_remove_defaults = arguments.length > 4 && arguments[4] !== void 0 ? arguments[4] : false, skip_warning = arguments.length > 5 && arguments[5] !== void 0 ? arguments[5] : false;
    if (hydrating && should_remove_defaults && element.tagName === 'INPUT') {
        var input = /** @type {HTMLInputElement} */ element;
        var attribute = input.type === 'checkbox' ? 'defaultChecked' : 'defaultValue';
        if (!(attribute in next)) {
            remove_input_defaults(input);
        }
    }
    var attributes = get_attributes(element);
    var is_custom_element = attributes[IS_CUSTOM_ELEMENT];
    var preserve_attribute_case = !attributes[IS_HTML];
    // If we're hydrating but the custom element is from Svelte, and it already scaffolded,
    // then it might run block logic in hydration mode, which we have to prevent.
    let is_hydrating_custom_element = hydrating && is_custom_element;
    if (is_hydrating_custom_element) {
        set_hydrating(false);
    }
    var current = prev || {};
    var is_option_element = element.tagName === 'OPTION';
    for(var key in prev){
        if (!(key in next)) {
            next[key] = null;
        }
    }
    if (next.class) {
        next.class = clsx(next.class);
    } else if (css_hash || next[CLASS]) {
        next.class = null; /* force call to set_class() */ 
    }
    if (next[STYLE]) {
        var _next;
        (_next = next).style ?? (_next.style = null); /* force call to set_style() */ 
    }
    var setters = get_setters(element);
    // since key is captured we use const
    for(const key in next){
        // let instead of var because referenced in a closure
        let value = next[key];
        // Up here because we want to do this for the initial value, too, even if it's undefined,
        // and this wouldn't be reached in case of undefined because of the equality check below
        if (is_option_element && key === 'value' && value == null) {
            // The <option> element is a special case because removing the value attribute means
            // the value is set to the text content of the option element, and setting the value
            // to null or undefined means the value is set to the string "null" or "undefined".
            // To align with how we handle this case in non-spread-scenarios, this logic is needed.
            // There's a super-edge-case bug here that is left in in favor of smaller code size:
            // Because of the "set missing props to null" logic above, we can't differentiate
            // between a missing value and an explicitly set value of null or undefined. That means
            // that once set, the value attribute of an <option> element can't be removed. This is
            // a very rare edge case, and removing the attribute altogether isn't possible either
            // for the <option value={undefined}> case, so we're not losing any functionality here.
            // @ts-ignore
            element.value = element.__value = '';
            current[key] = value;
            continue;
        }
        if (key === 'class') {
            var is_html = element.namespaceURI === 'http://www.w3.org/1999/xhtml';
            set_class(element, is_html, value, css_hash, prev === null || prev === void 0 ? void 0 : prev[CLASS], next[CLASS]);
            current[key] = value;
            current[CLASS] = next[CLASS];
            continue;
        }
        if (key === 'style') {
            set_style(element, value, prev === null || prev === void 0 ? void 0 : prev[STYLE], next[STYLE]);
            current[key] = value;
            current[STYLE] = next[STYLE];
            continue;
        }
        var prev_value = current[key];
        // Skip if value is unchanged, unless it's `undefined` and the element still has the attribute
        if (value === prev_value && !(value === undefined && element.hasAttribute(key))) {
            continue;
        }
        current[key] = value;
        var prefix = key[0] + key[1]; // this is faster than key.slice(0, 2)
        if (prefix === '$$') continue;
        if (prefix === 'on') {
            /** @type {{ capture?: true }} */ const opts = {};
            const event_handle_key = '$$' + key;
            let event_name = key.slice(2);
            var delegated = is_delegated(event_name);
            if (is_capture_event(event_name)) {
                event_name = event_name.slice(0, -7);
                opts.capture = true;
            }
            if (!delegated && prev_value) {
                // Listening to same event but different handler -> our handle function below takes care of this
                // If we were to remove and add listeners in this case, it could happen that the event is "swallowed"
                // (the browser seems to not know yet that a new one exists now) and doesn't reach the handler
                // https://github.com/sveltejs/svelte/issues/11903
                if (value != null) continue;
                element.removeEventListener(event_name, current[event_handle_key], opts);
                current[event_handle_key] = null;
            }
            if (value != null) {
                if (!delegated) {
                    /**
					 * @this {any}
					 * @param {Event} evt
					 */ function handle(evt) {
                        current[key].call(this, evt);
                    }
                    current[event_handle_key] = create_event(event_name, element, handle, opts);
                } else {
                    // @ts-ignore
                    element[`__${event_name}`] = value;
                    delegate([
                        event_name
                    ]);
                }
            } else if (delegated) {
                // @ts-ignore
                element[`__${event_name}`] = undefined;
            }
        } else if (key === 'style') {
            // avoid using the setter
            set_attribute(element, key, value);
        } else if (key === 'autofocus') {
            autofocus(/** @type {HTMLElement} */ element, Boolean(value));
        } else if (!is_custom_element && (key === '__value' || key === 'value' && value != null)) {
            // @ts-ignore We're not running this for custom elements because __value is actually
            // how Lit stores the current value on the element, and messing with that would break things.
            element.value = element.__value = value;
        } else if (key === 'selected' && is_option_element) {
            set_selected(/** @type {HTMLOptionElement} */ element, value);
        } else {
            var name = key;
            if (!preserve_attribute_case) {
                name = normalize_attribute(name);
            }
            var is_default = name === 'defaultValue' || name === 'defaultChecked';
            if (value == null && !is_custom_element && !is_default) {
                attributes[key] = null;
                if (name === 'value' || name === 'checked') {
                    // removing value/checked also removes defaultValue/defaultChecked — preserve
                    let input = /** @type {HTMLInputElement} */ element;
                    const use_default = prev === undefined;
                    if (name === 'value') {
                        let previous = input.defaultValue;
                        input.removeAttribute(name);
                        input.defaultValue = previous;
                        // @ts-ignore
                        input.value = input.__value = use_default ? previous : null;
                    } else {
                        let previous = input.defaultChecked;
                        input.removeAttribute(name);
                        input.defaultChecked = previous;
                        input.checked = use_default ? previous : false;
                    }
                } else {
                    element.removeAttribute(key);
                }
            } else if (is_default || setters.includes(name) && (is_custom_element || typeof value !== 'string')) {
                // @ts-ignore
                element[name] = value;
                // remove it from attributes's cache
                if (name in attributes) attributes[name] = UNINITIALIZED;
            } else if (typeof value !== 'function') {
                set_attribute(element, name, value, skip_warning);
            }
        }
    }
    if (is_hydrating_custom_element) {
        set_hydrating(true);
    }
    return current;
}
/**
 * @param {Element & ElementCSSInlineStyle} element
 * @param {(...expressions: any) => Record<string | symbol, any>} fn
 * @param {Array<() => any>} sync
 * @param {Array<() => Promise<any>>} async
 * @param {string} [css_hash]
 * @param {boolean} [should_remove_defaults]
 * @param {boolean} [skip_warning]
 */ function attribute_effect(element, fn) {
    let sync = arguments.length > 2 && arguments[2] !== void 0 ? arguments[2] : [], async = arguments.length > 3 && arguments[3] !== void 0 ? arguments[3] : [], css_hash = arguments.length > 4 ? arguments[4] : void 0, should_remove_defaults = arguments.length > 5 && arguments[5] !== void 0 ? arguments[5] : false, skip_warning = arguments.length > 6 && arguments[6] !== void 0 ? arguments[6] : false;
    flatten(sync, async, (values)=>{
        /** @type {Record<string | symbol, any> | undefined} */ var prev = undefined;
        /** @type {Record<symbol, Effect>} */ var effects = {};
        var is_select = element.nodeName === 'SELECT';
        var inited = false;
        block(()=>{
            var next = fn(...values.map(get));
            /** @type {Record<string | symbol, any>} */ var current = set_attributes(element, prev, next, css_hash, should_remove_defaults, skip_warning);
            if (inited && is_select && 'value' in next) {
                select_option(/** @type {HTMLSelectElement} */ element, next.value);
            }
            for (let symbol of Object.getOwnPropertySymbols(effects)){
                if (!next[symbol]) destroy_effect(effects[symbol]);
            }
            for (let symbol of Object.getOwnPropertySymbols(next)){
                var n = next[symbol];
                if (symbol.description === ATTACHMENT_KEY && (!prev || n !== prev[symbol])) {
                    if (effects[symbol]) destroy_effect(effects[symbol]);
                    effects[symbol] = branch(()=>attach(element, ()=>n));
                }
                current[symbol] = n;
            }
            prev = current;
        });
        if (is_select) {
            var select = /** @type {HTMLSelectElement} */ element;
            effect(()=>{
                select_option(select, /** @type {Record<string | symbol, any>} */ prev.value, true);
                init_select(select);
            });
        }
        inited = true;
    });
}
/**
 *
 * @param {Element} element
 */ function get_attributes(element) {
    var /** @type {Record<string | symbol, unknown>} **/ // @ts-expect-error
    _element;
    return (_element = element).__attributes ?? (_element.__attributes = {
        [IS_CUSTOM_ELEMENT]: element.nodeName.includes('-'),
        [IS_HTML]: element.namespaceURI === NAMESPACE_HTML
    });
}
/** @type {Map<string, string[]>} */ var setters_cache = new Map();
/** @param {Element} element */ function get_setters(element) {
    var cache_key = element.getAttribute('is') || element.nodeName;
    var setters = setters_cache.get(cache_key);
    if (setters) return setters;
    setters_cache.set(cache_key, setters = []);
    var descriptors;
    var proto = element; // In the case of custom elements there might be setters on the instance
    var element_proto = Element.prototype;
    // Stop at Element, from there on there's only unnecessary setters we're not interested in
    // Do not use contructor.name here as that's unreliable in some browser environments
    while(element_proto !== proto){
        descriptors = get_descriptors(proto);
        for(var key in descriptors){
            if (descriptors[key].set) {
                setters.push(key);
            }
        }
        proto = get_prototype_of(proto);
    }
    return setters;
}
/**
 * @param {any} element
 * @param {string} attribute
 * @param {string} value
 */ function check_src_in_dev_hydration(element, attribute, value) {
    if (!DEV) return;
    if (attribute === 'srcset' && srcset_url_equal(element, value)) return;
    if (src_url_equal(element.getAttribute(attribute) ?? '', value)) return;
    w.hydration_attribute_changed(attribute, element.outerHTML.replace(element.innerHTML, element.innerHTML && '...'), String(value));
}
/**
 * @param {string} element_src
 * @param {string} url
 * @returns {boolean}
 */ function src_url_equal(element_src, url) {
    if (element_src === url) return true;
    return new URL(element_src, document.baseURI).href === new URL(url, document.baseURI).href;
}
/** @param {string} srcset */ function split_srcset(srcset) {
    return srcset.split(',').map((src)=>src.trim().split(' ').filter(Boolean));
}
/**
 * @param {HTMLSourceElement | HTMLImageElement} element
 * @param {string} srcset
 * @returns {boolean}
 */ function srcset_url_equal(element, srcset) {
    var element_urls = split_srcset(element.srcset);
    var urls = split_srcset(srcset);
    return urls.length === element_urls.length && urls.every((param, i)=>{
        let [url, width] = param;
        return width === element_urls[i][1] && // We need to test both ways because Vite will create an a full URL with
        // `new URL(asset, import.meta.url).href` for the client when `base: './'`, and the
        // relative URLs inside srcset are not automatically resolved to absolute URLs by
        // browsers (in contrast to img.src). This means both SSR and DOM code could
        // contain relative or absolute URLs.
        (src_url_equal(element_urls[i][0], url) || src_url_equal(url, element_urls[i][0]));
    });
}

;// CONCATENATED MODULE: ./node_modules/.pnpm/esm-env@1.2.2/node_modules/esm-env/true.js
/* ESM default export */ const esm_env_true = (true);

;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/timing.js
/** @import { Raf } from '#client' */ 

const timing_now = esm_env_true ? ()=>performance.now() : ()=>Date.now();
/** @type {Raf} */ const timing_raf = {
    // don't access requestAnimationFrame eagerly outside method
    // this allows basic testing of user code without JSDOM
    // bunder will eval and remove ternary when the user's app is built
    tick: /** @param {any} _ */ (_)=>(esm_env_true ? requestAnimationFrame : shared_utils/* .noop */.lQ)(_),
    now: ()=>timing_now(),
    tasks: new Set()
};

;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/loop.js
/** @import { TaskCallback, Task, TaskEntry } from '#client' */ 
// TODO move this into timing.js where it probably belongs
/**
 * @returns {void}
 */ function run_tasks() {
    // use `raf.now()` instead of the `requestAnimationFrame` callback argument, because
    // otherwise things can get wonky https://github.com/sveltejs/svelte/pull/14541
    const now = raf.now();
    raf.tasks.forEach((task)=>{
        if (!task.c(now)) {
            raf.tasks.delete(task);
            task.f();
        }
    });
    if (raf.tasks.size !== 0) {
        raf.tick(run_tasks);
    }
}
/**
 * Creates a new task that runs on each raf frame
 * until it returns a falsy value or is aborted
 * @param {TaskCallback} callback
 * @returns {Task}
 */ function loop_loop(callback) {
    /** @type {TaskEntry} */ let task;
    if (raf.tasks.size === 0) {
        raf.tick(run_tasks);
    }
    return {
        promise: new Promise((fulfill)=>{
            raf.tasks.add(task = {
                c: callback,
                f: fulfill
            });
        }),
        abort () {
            raf.tasks.delete(task);
        }
    };
}

;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/elements/transitions.js
/** @import { AnimateFn, Animation, AnimationConfig, EachItem, Effect, TransitionFn, TransitionManager } from '#client' */ 









/**
 * @param {Element} element
 * @param {'introstart' | 'introend' | 'outrostart' | 'outroend'} type
 * @returns {void}
 */ function dispatch_event(element, type) {
    without_reactive_context(()=>{
        element.dispatchEvent(new CustomEvent(type));
    });
}
/**
 * Converts a property to the camel-case format expected by Element.animate(), KeyframeEffect(), and KeyframeEffect.setKeyframes().
 * @param {string} style
 * @returns {string}
 */ function css_property_to_camelcase(style) {
    // in compliance with spec
    if (style === 'float') return 'cssFloat';
    if (style === 'offset') return 'cssOffset';
    // do not rename custom @properties
    if (style.startsWith('--')) return style;
    const parts = style.split('-');
    if (parts.length === 1) return parts[0];
    return parts[0] + parts.slice(1).map(/** @param {any} word */ (word)=>word[0].toUpperCase() + word.slice(1)).join('');
}
/**
 * @param {string} css
 * @returns {Keyframe}
 */ function css_to_keyframe(css) {
    /** @type {Keyframe} */ const keyframe = {};
    const parts = css.split(';');
    for (const part of parts){
        const [property, value] = part.split(':');
        if (!property || value === undefined) break;
        const formatted_property = css_property_to_camelcase(property.trim());
        keyframe[formatted_property] = value.trim();
    }
    return keyframe;
}
/** @param {number} t */ const linear = (t)=>t;
/**
 * Called inside keyed `{#each ...}` blocks (as `$.animation(...)`). This creates an animation manager
 * and attaches it to the block, so that moves can be animated following reconciliation.
 * @template P
 * @param {Element} element
 * @param {() => AnimateFn<P | undefined>} get_fn
 * @param {(() => P) | null} get_params
 */ function transitions_animation(element, get_fn, get_params) {
    var _item;
    var item = /** @type {EachItem} */ current_each_item;
    /** @type {DOMRect} */ var from;
    /** @type {DOMRect} */ var to;
    /** @type {Animation | undefined} */ var animation;
    /** @type {null | { position: string, width: string, height: string, transform: string }} */ var original_styles = null;
    (_item = item).a ?? (_item.a = {
        element,
        measure () {
            from = this.element.getBoundingClientRect();
        },
        apply () {
            animation === null || animation === void 0 ? void 0 : animation.abort();
            to = this.element.getBoundingClientRect();
            if (from.left !== to.left || from.right !== to.right || from.top !== to.top || from.bottom !== to.bottom) {
                const options = get_fn()(this.element, {
                    from,
                    to
                }, get_params === null || get_params === void 0 ? void 0 : get_params());
                animation = animate(this.element, options, undefined, 1, ()=>{
                    animation === null || animation === void 0 ? void 0 : animation.abort();
                    animation = undefined;
                });
            }
        },
        fix () {
            // If an animation is already running, transforming the element is likely to fail,
            // because the styles applied by the animation take precedence. In the case of crossfade,
            // that means the `translate(...)` of the crossfade transition overrules the `translate(...)`
            // we would apply below, leading to the element jumping somewhere to the top left.
            if (element.getAnimations().length) return;
            // It's important to destructure these to get fixed values - the object itself has getters,
            // and changing the style to 'absolute' can for example influence the width.
            var { position, width, height } = getComputedStyle(element);
            if (position !== 'absolute' && position !== 'fixed') {
                var style = /** @type {HTMLElement | SVGElement} */ element.style;
                original_styles = {
                    position: style.position,
                    width: style.width,
                    height: style.height,
                    transform: style.transform
                };
                style.position = 'absolute';
                style.width = width;
                style.height = height;
                var to = element.getBoundingClientRect();
                if (from.left !== to.left || from.top !== to.top) {
                    var transform = `translate(${from.left - to.left}px, ${from.top - to.top}px)`;
                    style.transform = style.transform ? `${style.transform} ${transform}` : transform;
                }
            }
        },
        unfix () {
            if (original_styles) {
                var style = /** @type {HTMLElement | SVGElement} */ element.style;
                style.position = original_styles.position;
                style.width = original_styles.width;
                style.height = original_styles.height;
                style.transform = original_styles.transform;
            }
        }
    });
    // in the case of a `<svelte:element>`, it's possible for `$.animation(...)` to be called
    // when an animation manager already exists, if the tag changes. in that case, we need to
    // swap out the element rather than creating a new manager, in case it happened at the same
    // moment as a reconciliation
    item.a.element = element;
}
/**
 * Called inside block effects as `$.transition(...)`. This creates a transition manager and
 * attaches it to the current effect — later, inside `pause_effect` and `resume_effect`, we
 * use this to create `intro` and `outro` transitions.
 * @template P
 * @param {number} flags
 * @param {HTMLElement} element
 * @param {() => TransitionFn<P | undefined>} get_fn
 * @param {(() => P) | null} get_params
 * @returns {void}
 */ function transitions_transition(flags, element, get_fn, get_params) {
    var _e;
    var is_intro = (flags & TRANSITION_IN) !== 0;
    var is_outro = (flags & TRANSITION_OUT) !== 0;
    var is_both = is_intro && is_outro;
    var is_global = (flags & TRANSITION_GLOBAL) !== 0;
    /** @type {'in' | 'out' | 'both'} */ var direction = is_both ? 'both' : is_intro ? 'in' : 'out';
    /** @type {AnimationConfig | ((opts: { direction: 'in' | 'out' }) => AnimationConfig) | undefined} */ var current_options;
    var inert = element.inert;
    /**
	 * The default overflow style, stashed so we can revert changes during the transition
	 * that are necessary to work around a Safari <18 bug
	 * TODO 6.0 remove this, if older versions of Safari have died out enough
	 */ var overflow = element.style.overflow;
    /** @type {Animation | undefined} */ var intro;
    /** @type {Animation | undefined} */ var outro;
    function get_options() {
        return without_reactive_context(()=>{
            // If a transition is still ongoing, we use the existing options rather than generating
            // new ones. This ensures that reversible transitions reverse smoothly, rather than
            // jumping to a new spot because (for example) a different `duration` was used
            return current_options ?? (current_options = get_fn()(element, (get_params === null || get_params === void 0 ? void 0 : get_params()) ?? /** @type {P} */ {}, {
                direction
            }));
        });
    }
    /** @type {TransitionManager} */ var transition = {
        is_global,
        in () {
            element.inert = inert;
            if (!is_intro) {
                var _outro_reset;
                outro === null || outro === void 0 ? void 0 : outro.abort();
                outro === null || outro === void 0 ? void 0 : (_outro_reset = outro.reset) === null || _outro_reset === void 0 ? void 0 : _outro_reset.call(outro);
                return;
            }
            if (!is_outro) {
                // if we intro then outro then intro again, we want to abort the first intro,
                // if it's not a bidirectional transition
                intro === null || intro === void 0 ? void 0 : intro.abort();
            }
            dispatch_event(element, 'introstart');
            intro = animate(element, get_options(), outro, 1, ()=>{
                dispatch_event(element, 'introend');
                // Ensure we cancel the animation to prevent leaking
                intro === null || intro === void 0 ? void 0 : intro.abort();
                intro = current_options = undefined;
                element.style.overflow = overflow;
            });
        },
        out (fn) {
            if (!is_outro) {
                fn === null || fn === void 0 ? void 0 : fn();
                current_options = undefined;
                return;
            }
            element.inert = true;
            dispatch_event(element, 'outrostart');
            outro = animate(element, get_options(), intro, 0, ()=>{
                dispatch_event(element, 'outroend');
                fn === null || fn === void 0 ? void 0 : fn();
            });
        },
        stop: ()=>{
            intro === null || intro === void 0 ? void 0 : intro.abort();
            outro === null || outro === void 0 ? void 0 : outro.abort();
        }
    };
    var e = /** @type {Effect} */ active_effect;
    ((_e = e).transitions ?? (_e.transitions = [])).push(transition);
    // if this is a local transition, we only want to run it if the parent (branch) effect's
    // parent (block) effect is where the state change happened. we can determine that by
    // looking at whether the block effect is currently initializing
    if (is_intro && should_intro) {
        var run = is_global;
        if (!run) {
            var block = /** @type {Effect | null} */ e.parent;
            // skip over transparent blocks (e.g. snippets, else-if blocks)
            while(block && (block.f & EFFECT_TRANSPARENT) !== 0){
                while(block = block.parent){
                    if ((block.f & BLOCK_EFFECT) !== 0) break;
                }
            }
            run = !block || (block.f & EFFECT_RAN) !== 0;
        }
        if (run) {
            effect(()=>{
                untrack(()=>transition.in());
            });
        }
    }
}
/**
 * Animates an element, according to the provided configuration
 * @param {Element} element
 * @param {AnimationConfig | ((opts: { direction: 'in' | 'out' }) => AnimationConfig)} options
 * @param {Animation | undefined} counterpart The corresponding intro/outro to this outro/intro
 * @param {number} t2 The target `t` value — `1` for intro, `0` for outro
 * @param {(() => void)} on_finish Called after successfully completing the animation
 * @returns {Animation}
 */ function animate(element, options, counterpart, t2, on_finish) {
    var is_intro = t2 === 1;
    if (is_function(options)) {
        // In the case of a deferred transition (such as `crossfade`), `option` will be
        // a function rather than an `AnimationConfig`. We need to call this function
        // once the DOM has been updated...
        /** @type {Animation} */ var a;
        var aborted = false;
        queue_micro_task(()=>{
            if (aborted) return;
            var o = options({
                direction: is_intro ? 'in' : 'out'
            });
            a = animate(element, o, counterpart, t2, on_finish);
        });
        // ...but we want to do so without using `async`/`await` everywhere, so
        // we return a facade that allows everything to remain synchronous
        return {
            abort: ()=>{
                aborted = true;
                a === null || a === void 0 ? void 0 : a.abort();
            },
            deactivate: ()=>a.deactivate(),
            reset: ()=>a.reset(),
            t: ()=>a.t()
        };
    }
    counterpart === null || counterpart === void 0 ? void 0 : counterpart.deactivate();
    if (!(options === null || options === void 0 ? void 0 : options.duration)) {
        on_finish();
        return {
            abort: noop,
            deactivate: noop,
            reset: noop,
            t: ()=>t2
        };
    }
    const { delay = 0, css, tick, easing = linear } = options;
    var keyframes = [];
    if (is_intro && counterpart === undefined) {
        if (tick) {
            tick(0, 1); // TODO put in nested effect, to avoid interleaved reads/writes?
        }
        if (css) {
            var styles = css_to_keyframe(css(0, 1));
            keyframes.push(styles, styles);
        }
    }
    var get_t = ()=>1 - t2;
    // create a dummy animation that lasts as long as the delay (but with whatever devtools
    // multiplier is in effect). in the common case that it is `0`, we keep it anyway so that
    // the CSS keyframes aren't created until the DOM is updated
    //
    // fill forwards to prevent the element from rendering without styles applied
    // see https://github.com/sveltejs/svelte/issues/14732
    var animation = element.animate(keyframes, {
        duration: delay,
        fill: 'forwards'
    });
    animation.onfinish = ()=>{
        // remove dummy animation from the stack to prevent conflict with main animation
        animation.cancel();
        // for bidirectional transitions, we start from the current position,
        // rather than doing a full intro/outro
        var t1 = (counterpart === null || counterpart === void 0 ? void 0 : counterpart.t()) ?? 1 - t2;
        counterpart === null || counterpart === void 0 ? void 0 : counterpart.abort();
        var delta = t2 - t1;
        var duration = /** @type {number} */ options.duration * Math.abs(delta);
        var keyframes = [];
        if (duration > 0) {
            /**
			 * Whether or not the CSS includes `overflow: hidden`, in which case we need to
			 * add it as an inline style to work around a Safari <18 bug
			 * TODO 6.0 remove this, if possible
			 */ var needs_overflow_hidden = false;
            if (css) {
                var n = Math.ceil(duration / (1000 / 60)); // `n` must be an integer, or we risk missing the `t2` value
                for(var i = 0; i <= n; i += 1){
                    var t = t1 + delta * easing(i / n);
                    var styles = css_to_keyframe(css(t, 1 - t));
                    keyframes.push(styles);
                    needs_overflow_hidden || (needs_overflow_hidden = styles.overflow === 'hidden');
                }
            }
            if (needs_overflow_hidden) {
                /** @type {HTMLElement} */ element.style.overflow = 'hidden';
            }
            get_t = ()=>{
                var time = /** @type {number} */ /** @type {globalThis.Animation} */ animation.currentTime;
                return t1 + delta * easing(time / duration);
            };
            if (tick) {
                loop(()=>{
                    if (animation.playState !== 'running') return false;
                    var t = get_t();
                    tick(t, 1 - t);
                    return true;
                });
            }
        }
        animation = element.animate(keyframes, {
            duration,
            fill: 'forwards'
        });
        animation.onfinish = ()=>{
            get_t = ()=>t2;
            tick === null || tick === void 0 ? void 0 : tick(t2, 1 - t2);
            on_finish();
        };
    };
    return {
        abort: ()=>{
            if (animation) {
                animation.cancel();
                // This prevents memory leaks in Chromium
                animation.effect = null;
                // This prevents onfinish to be launched after cancel(),
                // which can happen in some rare cases
                // see https://github.com/sveltejs/svelte/issues/13681
                animation.onfinish = noop;
            }
        },
        deactivate: ()=>{
            on_finish = noop;
        },
        reset: ()=>{
            if (t2 === 0) {
                tick === null || tick === void 0 ? void 0 : tick(1, 0);
            }
        },
        t: ()=>get_t()
    };
}

;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/elements/bindings/document.js

/**
 * @param {(activeElement: Element | null) => void} update
 * @returns {void}
 */ function bind_active_element(update) {
    listen(document, [
        'focusin',
        'focusout'
    ], (event)=>{
        if (event && event.type === 'focusout' && /** @type {FocusEvent} */ event.relatedTarget) {
            // The tests still pass if we remove this, because of JSDOM limitations, but it is necessary
            // to avoid temporarily resetting to `document.body`
            return;
        }
        update(document.activeElement);
    });
}

;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/elements/bindings/input.js
/** @import { Batch } from '../../../reactivity/batch.js' */ 








/**
 * @param {HTMLInputElement} input
 * @param {() => unknown} get
 * @param {(value: unknown) => void} set
 * @returns {void}
 */ function bind_value(input, get) {
    let set = arguments.length > 2 && arguments[2] !== void 0 ? arguments[2] : get;
    var batches = new WeakSet();
    listen_to_event_and_reset_event(input, 'input', async (is_reset)=>{
        if (DEV && input.type === 'checkbox') {
            // TODO should this happen in prod too?
            e.bind_invalid_checkbox_value();
        }
        /** @type {any} */ var value = is_reset ? input.defaultValue : input.value;
        value = is_numberlike_input(input) ? to_number(value) : value;
        set(value);
        if (current_batch !== null) {
            batches.add(current_batch);
        }
        // Because `{#each ...}` blocks work by updating sources inside the flush,
        // we need to wait a tick before checking to see if we should forcibly
        // update the input and reset the selection state
        await tick();
        // Respect any validation in accessors
        if (value !== (value = get())) {
            var start = input.selectionStart;
            var end = input.selectionEnd;
            // the value is coerced on assignment
            input.value = value ?? '';
            // Restore selection
            if (end !== null) {
                input.selectionStart = start;
                input.selectionEnd = Math.min(end, input.value.length);
            }
        }
    });
    if (// If we are hydrating and the value has since changed,
    // then use the updated value from the input instead.
    hydrating && input.defaultValue !== input.value || // If defaultValue is set, then value == defaultValue
    // TODO Svelte 6: remove input.value check and set to empty string?
    untrack(get) == null && input.value) {
        set(is_numberlike_input(input) ? to_number(input.value) : input.value);
        if (current_batch !== null) {
            batches.add(current_batch);
        }
    }
    render_effect(()=>{
        if (DEV && input.type === 'checkbox') {
            // TODO should this happen in prod too?
            e.bind_invalid_checkbox_value();
        }
        var value = get();
        if (input === document.activeElement) {
            // we need both, because in non-async mode, render effects run before previous_batch is set
            var batch = /** @type {Batch} */ previous_batch ?? current_batch;
            // Never rewrite the contents of a focused input. We can get here if, for example,
            // an update is deferred because of async work depending on the input:
            //
            // <input bind:value={query}>
            // <p>{await find(query)}</p>
            if (batches.has(batch)) {
                return;
            }
        }
        if (is_numberlike_input(input) && value === to_number(input.value)) {
            // handles 0 vs 00 case (see https://github.com/sveltejs/svelte/issues/9959)
            return;
        }
        if (input.type === 'date' && !value && !input.value) {
            // Handles the case where a temporarily invalid date is set (while typing, for example with a leading 0 for the day)
            // and prevents this state from clearing the other parts of the date input (see https://github.com/sveltejs/svelte/issues/7897)
            return;
        }
        // don't set the value of the input if it's the same to allow
        // minlength to work properly
        if (value !== input.value) {
            // @ts-expect-error the value is coerced on assignment
            input.value = value ?? '';
        }
    });
}
/** @type {Set<HTMLInputElement[]>} */ const input_pending = new Set();
/**
 * @param {HTMLInputElement[]} inputs
 * @param {null | [number]} group_index
 * @param {HTMLInputElement} input
 * @param {() => unknown} get
 * @param {(value: unknown) => void} set
 * @returns {void}
 */ function bind_group(inputs, group_index, input, get) {
    let set = arguments.length > 4 && arguments[4] !== void 0 ? arguments[4] : get;
    var is_checkbox = input.getAttribute('type') === 'checkbox';
    var binding_group = inputs;
    // needs to be let or related code isn't treeshaken out if it's always false
    let hydration_mismatch = false;
    if (group_index !== null) {
        for (var index of group_index){
            var _binding_group, _index;
            // @ts-expect-error
            binding_group = (_binding_group = binding_group)[_index = index] ?? (_binding_group[_index] = []);
        }
    }
    binding_group.push(input);
    listen_to_event_and_reset_event(input, 'change', ()=>{
        // @ts-ignore
        var value = input.__value;
        if (is_checkbox) {
            value = get_binding_group_value(binding_group, value, input.checked);
        }
        set(value);
    }, // TODO better default value handling
    ()=>set(is_checkbox ? [] : null));
    render_effect(()=>{
        var value = get();
        // If we are hydrating and the value has since changed, then use the update value
        // from the input instead.
        if (hydrating && input.defaultChecked !== input.checked) {
            hydration_mismatch = true;
            return;
        }
        if (is_checkbox) {
            value = value || [];
            // @ts-ignore
            input.checked = value.includes(input.__value);
        } else {
            // @ts-ignore
            input.checked = is(input.__value, value);
        }
    });
    teardown(()=>{
        var index = binding_group.indexOf(input);
        if (index !== -1) {
            binding_group.splice(index, 1);
        }
    });
    if (!input_pending.has(binding_group)) {
        input_pending.add(binding_group);
        queue_micro_task(()=>{
            // necessary to maintain binding group order in all insertion scenarios
            binding_group.sort((a, b)=>a.compareDocumentPosition(b) === 4 ? -1 : 1);
            input_pending.delete(binding_group);
        });
    }
    queue_micro_task(()=>{
        if (hydration_mismatch) {
            var value;
            if (is_checkbox) {
                value = get_binding_group_value(binding_group, value, input.checked);
            } else {
                var hydration_input = binding_group.find((input)=>input.checked);
                // @ts-ignore
                value = hydration_input === null || hydration_input === void 0 ? void 0 : hydration_input.__value;
            }
            set(value);
        }
    });
}
/**
 * @param {HTMLInputElement} input
 * @param {() => unknown} get
 * @param {(value: unknown) => void} set
 * @returns {void}
 */ function bind_checked(input, get) {
    let set = arguments.length > 2 && arguments[2] !== void 0 ? arguments[2] : get;
    listen_to_event_and_reset_event(input, 'change', (is_reset)=>{
        var value = is_reset ? input.defaultChecked : input.checked;
        set(value);
    });
    if (// If we are hydrating and the value has since changed,
    // then use the update value from the input instead.
    hydrating && input.defaultChecked !== input.checked || // If defaultChecked is set, then checked == defaultChecked
    untrack(get) == null) {
        set(input.checked);
    }
    render_effect(()=>{
        var value = get();
        input.checked = Boolean(value);
    });
}
/**
 * @template V
 * @param {Array<HTMLInputElement>} group
 * @param {V} __value
 * @param {boolean} checked
 * @returns {V[]}
 */ function get_binding_group_value(group, __value, checked) {
    /** @type {Set<V>} */ var value = new Set();
    for(var i = 0; i < group.length; i += 1){
        if (group[i].checked) {
            // @ts-ignore
            value.add(group[i].__value);
        }
    }
    if (!checked) {
        value.delete(__value);
    }
    return Array.from(value);
}
/**
 * @param {HTMLInputElement} input
 */ function is_numberlike_input(input) {
    var type = input.type;
    return type === 'number' || type === 'range';
}
/**
 * @param {string} value
 */ function to_number(value) {
    return value === '' ? null : +value;
}
/**
 * @param {HTMLInputElement} input
 * @param {() => FileList | null} get
 * @param {(value: FileList | null) => void} set
 */ function bind_files(input, get) {
    let set = arguments.length > 2 && arguments[2] !== void 0 ? arguments[2] : get;
    listen_to_event_and_reset_event(input, 'change', ()=>{
        set(input.files);
    });
    if (// If we are hydrating and the value has since changed,
    // then use the updated value from the input instead.
    hydrating && input.files) {
        set(input.files);
    }
    render_effect(()=>{
        input.files = get();
    });
}

;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/elements/bindings/media.js


/** @param {TimeRanges} ranges */ function time_ranges_to_array(ranges) {
    var array = [];
    for(var i = 0; i < ranges.length; i += 1){
        array.push({
            start: ranges.start(i),
            end: ranges.end(i)
        });
    }
    return array;
}
/**
 * @param {HTMLVideoElement | HTMLAudioElement} media
 * @param {() => number | undefined} get
 * @param {(value: number) => void} set
 * @returns {void}
 */ function bind_current_time(media, get) {
    let set = arguments.length > 2 && arguments[2] !== void 0 ? arguments[2] : get;
    /** @type {number} */ var raf_id;
    /** @type {number} */ var value;
    // Ideally, listening to timeupdate would be enough, but it fires too infrequently for the currentTime
    // binding, which is why we use a raf loop, too. We additionally still listen to timeupdate because
    // the user could be scrubbing through the video using the native controls when the media is paused.
    var callback = ()=>{
        cancelAnimationFrame(raf_id);
        if (!media.paused) {
            raf_id = requestAnimationFrame(callback);
        }
        var next_value = media.currentTime;
        if (value !== next_value) {
            set(value = next_value);
        }
    };
    raf_id = requestAnimationFrame(callback);
    media.addEventListener('timeupdate', callback);
    render_effect(()=>{
        var next_value = Number(get());
        if (value !== next_value && !isNaN(/** @type {any} */ next_value)) {
            media.currentTime = value = next_value;
        }
    });
    teardown(()=>{
        cancelAnimationFrame(raf_id);
        media.removeEventListener('timeupdate', callback);
    });
}
/**
 * @param {HTMLVideoElement | HTMLAudioElement} media
 * @param {(array: Array<{ start: number; end: number }>) => void} set
 */ function bind_buffered(media, set) {
    /** @type {{ start: number; end: number; }[]} */ var current;
    // `buffered` can update without emitting any event, so we check it on various events.
    // By specs, `buffered` always returns a new object, so we have to compare deeply.
    listen(media, [
        'loadedmetadata',
        'progress',
        'timeupdate',
        'seeking'
    ], ()=>{
        var ranges = media.buffered;
        if (!current || current.length !== ranges.length || current.some((range, i)=>ranges.start(i) !== range.start || ranges.end(i) !== range.end)) {
            current = time_ranges_to_array(ranges);
            set(current);
        }
    });
}
/**
 * @param {HTMLVideoElement | HTMLAudioElement} media
 * @param {(array: Array<{ start: number; end: number }>) => void} set
 */ function bind_seekable(media, set) {
    listen(media, [
        'loadedmetadata'
    ], ()=>set(time_ranges_to_array(media.seekable)));
}
/**
 * @param {HTMLVideoElement | HTMLAudioElement} media
 * @param {(array: Array<{ start: number; end: number }>) => void} set
 */ function bind_played(media, set) {
    listen(media, [
        'timeupdate'
    ], ()=>set(time_ranges_to_array(media.played)));
}
/**
 * @param {HTMLVideoElement | HTMLAudioElement} media
 * @param {(seeking: boolean) => void} set
 */ function bind_seeking(media, set) {
    listen(media, [
        'seeking',
        'seeked'
    ], ()=>set(media.seeking));
}
/**
 * @param {HTMLVideoElement | HTMLAudioElement} media
 * @param {(seeking: boolean) => void} set
 */ function bind_ended(media, set) {
    listen(media, [
        'timeupdate',
        'ended'
    ], ()=>set(media.ended));
}
/**
 * @param {HTMLVideoElement | HTMLAudioElement} media
 * @param {(ready_state: number) => void} set
 */ function bind_ready_state(media, set) {
    listen(media, [
        'loadedmetadata',
        'loadeddata',
        'canplay',
        'canplaythrough',
        'playing',
        'waiting',
        'emptied'
    ], ()=>set(media.readyState));
}
/**
 * @param {HTMLVideoElement | HTMLAudioElement} media
 * @param {() => number | undefined} get
 * @param {(playback_rate: number) => void} set
 */ function bind_playback_rate(media, get) {
    let set = arguments.length > 2 && arguments[2] !== void 0 ? arguments[2] : get;
    // Needs to happen after element is inserted into the dom (which is guaranteed by using effect),
    // else playback will be set back to 1 by the browser
    effect(()=>{
        var value = Number(get());
        if (value !== media.playbackRate && !isNaN(value)) {
            media.playbackRate = value;
        }
    });
    // Start listening to ratechange events after the element is inserted into the dom,
    // else playback will be set to 1 by the browser
    effect(()=>{
        listen(media, [
            'ratechange'
        ], ()=>{
            set(media.playbackRate);
        });
    });
}
/**
 * @param {HTMLVideoElement | HTMLAudioElement} media
 * @param {() => boolean | undefined} get
 * @param {(paused: boolean) => void} set
 */ function bind_paused(media, get) {
    let set = arguments.length > 2 && arguments[2] !== void 0 ? arguments[2] : get;
    var paused = get();
    var update = ()=>{
        if (paused !== media.paused) {
            set(paused = media.paused);
        }
    };
    // If someone switches the src while media is playing, the player will pause.
    // Listen to the canplay event to get notified of this situation.
    listen(media, [
        'play',
        'pause',
        'canplay'
    ], update, paused == null);
    // Needs to be an effect to ensure media element is mounted: else, if paused is `false` (i.e. should play right away)
    // a "The play() request was interrupted by a new load request" error would be thrown because the resource isn't loaded yet.
    effect(()=>{
        if ((paused = !!get()) !== media.paused) {
            if (paused) {
                media.pause();
            } else {
                media.play().catch(()=>{
                    set(paused = true);
                });
            }
        }
    });
}
/**
 * @param {HTMLVideoElement | HTMLAudioElement} media
 * @param {() => number | undefined} get
 * @param {(volume: number) => void} set
 */ function bind_volume(media, get) {
    let set = arguments.length > 2 && arguments[2] !== void 0 ? arguments[2] : get;
    var callback = ()=>{
        set(media.volume);
    };
    if (get() == null) {
        callback();
    }
    listen(media, [
        'volumechange'
    ], callback, false);
    render_effect(()=>{
        var value = Number(get());
        if (value !== media.volume && !isNaN(value)) {
            media.volume = value;
        }
    });
}
/**
 * @param {HTMLVideoElement | HTMLAudioElement} media
 * @param {() => boolean | undefined} get
 * @param {(muted: boolean) => void} set
 */ function bind_muted(media, get) {
    let set = arguments.length > 2 && arguments[2] !== void 0 ? arguments[2] : get;
    var callback = ()=>{
        set(media.muted);
    };
    if (get() == null) {
        callback();
    }
    listen(media, [
        'volumechange'
    ], callback, false);
    render_effect(()=>{
        var value = !!get();
        if (media.muted !== value) media.muted = value;
    });
}

;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/elements/bindings/navigator.js

/**
 * @param {(online: boolean) => void} update
 * @returns {void}
 */ function bind_online(update) {
    listen(window, [
        'online',
        'offline'
    ], ()=>{
        update(navigator.onLine);
    });
}

;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/elements/bindings/props.js


/**
 * Makes an `export`ed (non-prop) variable available on the `$$props` object
 * so that consumers can do `bind:x` on the component.
 * @template V
 * @param {Record<string, unknown>} props
 * @param {string} prop
 * @param {V} value
 * @returns {void}
 */ function bind_prop(props, prop, value) {
    var desc = get_descriptor(props, prop);
    if (desc && desc.set) {
        props[prop] = value;
        teardown(()=>{
            props[prop] = null;
        });
    }
}

// EXTERNAL MODULE: ./node_modules/.pnpm/@swc+helpers@0.5.17/node_modules/@swc/helpers/esm/_class_private_field_get.js + 1 modules
var _class_private_field_get = __webpack_require__(570);
// EXTERNAL MODULE: ./node_modules/.pnpm/@swc+helpers@0.5.17/node_modules/@swc/helpers/esm/_class_private_field_init.js
var _class_private_field_init = __webpack_require__(636);
// EXTERNAL MODULE: ./node_modules/.pnpm/@swc+helpers@0.5.17/node_modules/@swc/helpers/esm/_class_private_field_set.js + 1 modules
var _class_private_field_set = __webpack_require__(549);
// EXTERNAL MODULE: ./node_modules/.pnpm/@swc+helpers@0.5.17/node_modules/@swc/helpers/esm/_class_private_method_get.js
var _class_private_method_get = __webpack_require__(585);
// EXTERNAL MODULE: ./node_modules/.pnpm/@swc+helpers@0.5.17/node_modules/@swc/helpers/esm/_class_private_method_init.js
var _class_private_method_init = __webpack_require__(23);
// EXTERNAL MODULE: ./node_modules/.pnpm/@swc+helpers@0.5.17/node_modules/@swc/helpers/esm/_define_property.js
var _define_property = __webpack_require__(925);
;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/elements/bindings/size.js








var /** */ _listeners = /*#__PURE__*/ new WeakMap(), /** @type {ResizeObserver | undefined} */ _observer = /*#__PURE__*/ new WeakMap(), /** @type {ResizeObserverOptions} */ _options = /*#__PURE__*/ new WeakMap(), _getObserver = /*#__PURE__*/ new WeakSet();
/**
 * Resize observer singleton.
 * One listener per element only!
 * https://groups.google.com/a/chromium.org/g/blink-dev/c/z6ienONUb5A/m/F5-VcUZtBAAJ
 */ class ResizeObserverSingleton {
    /**
	 * @param {Element} element
	 * @param {(entry: ResizeObserverEntry) => any} listener
	 */ observe(element, listener) {
        var listeners = (0,_class_private_field_get._)(this, _listeners).get(element) || new Set();
        listeners.add(listener);
        (0,_class_private_field_get._)(this, _listeners).set(element, listeners);
        (0,_class_private_method_get._)(this, _getObserver, getObserver).call(this).observe(element, (0,_class_private_field_get._)(this, _options));
        return ()=>{
            var listeners = (0,_class_private_field_get._)(this, _listeners).get(element);
            listeners.delete(listener);
            if (listeners.size === 0) {
                (0,_class_private_field_get._)(this, _listeners).delete(element);
                (0,_class_private_field_get._)(/** @type {ResizeObserver} */ this, _observer).unobserve(element);
            }
        };
    }
    /** @param {ResizeObserverOptions} options */ constructor(options){
        (0,_class_private_method_init._)(this, _getObserver);
        (0,_class_private_field_init._)(this, _listeners, {
            writable: true,
            value: new WeakMap()
        });
        (0,_class_private_field_init._)(this, _observer, {
            writable: true,
            value: void 0
        });
        (0,_class_private_field_init._)(this, _options, {
            writable: true,
            value: void 0
        });
        (0,_class_private_field_set._)(this, _options, options);
    }
}
/** @static */ (0,_define_property._)(ResizeObserverSingleton, "entries", new WeakMap());
function getObserver() {
    return (0,_class_private_field_get._)(this, _observer) ?? (0,_class_private_field_set._)(this, _observer, new ResizeObserver(/** @param {any} entries */ (entries)=>{
        for (var entry of entries){
            ResizeObserverSingleton.entries.set(entry.target, entry);
            for (var listener of (0,_class_private_field_get._)(this, _listeners).get(entry.target) || []){
                listener(entry);
            }
        }
    }));
}
var resize_observer_content_box = /* @__PURE__ */ new ResizeObserverSingleton({
    box: 'content-box'
});
var resize_observer_border_box = /* @__PURE__ */ new ResizeObserverSingleton({
    box: 'border-box'
});
var resize_observer_device_pixel_content_box = /* @__PURE__ */ new ResizeObserverSingleton({
    box: 'device-pixel-content-box'
});
/**
 * @param {Element} element
 * @param {'contentRect' | 'contentBoxSize' | 'borderBoxSize' | 'devicePixelContentBoxSize'} type
 * @param {(entry: keyof ResizeObserverEntry) => void} set
 */ function bind_resize_observer(element, type, set) {
    var observer = type === 'contentRect' || type === 'contentBoxSize' ? resize_observer_content_box : type === 'borderBoxSize' ? resize_observer_border_box : resize_observer_device_pixel_content_box;
    var unsub = observer.observe(element, /** @param {any} entry */ (entry)=>set(entry[type]));
    teardown(unsub);
}
/**
 * @param {HTMLElement} element
 * @param {'clientWidth' | 'clientHeight' | 'offsetWidth' | 'offsetHeight'} type
 * @param {(size: number) => void} set
 */ function bind_element_size(element, type, set) {
    var unsub = resize_observer_border_box.observe(element, ()=>set(element[type]));
    effect(()=>{
        // The update could contain reads which should be ignored
        untrack(()=>set(element[type]));
        return unsub;
    });
}

;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/elements/bindings/this.js




/**
 * @param {any} bound_value
 * @param {Element} element_or_component
 * @returns {boolean}
 */ function is_bound_this(bound_value, element_or_component) {
    return bound_value === element_or_component || (bound_value === null || bound_value === void 0 ? void 0 : bound_value[STATE_SYMBOL]) === element_or_component;
}
/**
 * @param {any} element_or_component
 * @param {(value: unknown, ...parts: unknown[]) => void} update
 * @param {(...parts: unknown[]) => unknown} get_value
 * @param {() => unknown[]} [get_parts] Set if the this binding is used inside an each block,
 * 										returns all the parts of the each block context that are used in the expression
 * @returns {void}
 */ function bind_this() {
    let element_or_component = arguments.length > 0 && arguments[0] !== void 0 ? arguments[0] : {}, update = arguments.length > 1 ? arguments[1] : void 0, get_value = arguments.length > 2 ? arguments[2] : void 0, get_parts = arguments.length > 3 ? arguments[3] : void 0;
    effect(()=>{
        /** @type {unknown[]} */ var old_parts;
        /** @type {unknown[]} */ var parts;
        render_effect(()=>{
            old_parts = parts;
            // We only track changes to the parts, not the value itself to avoid unnecessary reruns.
            parts = (get_parts === null || get_parts === void 0 ? void 0 : get_parts()) || [];
            untrack(()=>{
                if (element_or_component !== get_value(...parts)) {
                    update(element_or_component, ...parts);
                    // If this is an effect rerun (cause: each block context changes), then nullfiy the binding at
                    // the previous position if it isn't already taken over by a different effect.
                    if (old_parts && is_bound_this(get_value(...old_parts), element_or_component)) {
                        update(null, ...old_parts);
                    }
                }
            });
        });
        return ()=>{
            // We cannot use effects in the teardown phase, we we use a microtask instead.
            queue_micro_task(()=>{
                if (parts && is_bound_this(get_value(...parts), element_or_component)) {
                    update(null, ...parts);
                }
            });
        };
    });
    return element_or_component;
}

;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/elements/bindings/universal.js


/**
 * @param {'innerHTML' | 'textContent' | 'innerText'} property
 * @param {HTMLElement} element
 * @param {() => unknown} get
 * @param {(value: unknown) => void} set
 * @returns {void}
 */ function bind_content_editable(property, element, get) {
    let set = arguments.length > 3 && arguments[3] !== void 0 ? arguments[3] : get;
    element.addEventListener('input', ()=>{
        // @ts-ignore
        set(element[property]);
    });
    render_effect(()=>{
        var value = get();
        if (element[property] !== value) {
            if (value == null) {
                // @ts-ignore
                var non_null_value = element[property];
                set(non_null_value);
            } else {
                // @ts-ignore
                element[property] = value + '';
            }
        }
    });
}
/**
 * @param {string} property
 * @param {string} event_name
 * @param {Element} element
 * @param {(value: unknown) => void} set
 * @param {() => unknown} [get]
 * @returns {void}
 */ function bind_property(property, event_name, element, set, get) {
    var handler = ()=>{
        // @ts-ignore
        set(element[property]);
    };
    element.addEventListener(event_name, handler);
    if (get) {
        render_effect(()=>{
            // @ts-ignore
            element[property] = get();
        });
    } else {
        handler();
    }
    // @ts-ignore
    if (element === document.body || element === window || element === document) {
        teardown(()=>{
            element.removeEventListener(event_name, handler);
        });
    }
}
/**
 * @param {HTMLElement} element
 * @param {(value: unknown) => void} set
 * @returns {void}
 */ function bind_focused(element, set) {
    listen(element, [
        'focus',
        'blur'
    ], ()=>{
        set(element === document.activeElement);
    });
}

;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/elements/bindings/window.js


/**
 * @param {'x' | 'y'} type
 * @param {() => number} get
 * @param {(value: number) => void} set
 * @returns {void}
 */ function bind_window_scroll(type, get) {
    let set = arguments.length > 2 && arguments[2] !== void 0 ? arguments[2] : get;
    var is_scrolling_x = type === 'x';
    var target_handler = ()=>without_reactive_context(()=>{
            scrolling = true;
            clearTimeout(timeout);
            timeout = setTimeout(clear, 100); // TODO use scrollend event if supported (or when supported everywhere?)
            set(window[is_scrolling_x ? 'scrollX' : 'scrollY']);
        });
    addEventListener('scroll', target_handler, {
        passive: true
    });
    var scrolling = false;
    /** @type {ReturnType<typeof setTimeout>} */ var timeout;
    var clear = ()=>{
        scrolling = false;
    };
    var first = true;
    render_effect(()=>{
        var latest_value = get();
        // Don't scroll to the initial value for accessibility reasons
        if (first) {
            first = false;
        } else if (!scrolling && latest_value != null) {
            scrolling = true;
            clearTimeout(timeout);
            if (is_scrolling_x) {
                scrollTo(latest_value, window.scrollY);
            } else {
                scrollTo(window.scrollX, latest_value);
            }
            timeout = setTimeout(clear, 100);
        }
    });
    // Browsers don't fire the scroll event for the initial scroll position when scroll style isn't set to smooth
    effect(target_handler);
    teardown(()=>{
        removeEventListener('scroll', target_handler);
    });
}
/**
 * @param {'innerWidth' | 'innerHeight' | 'outerWidth' | 'outerHeight'} type
 * @param {(size: number) => void} set
 */ function bind_window_size(type, set) {
    listen(window, [
        'resize'
    ], ()=>without_reactive_context(()=>set(window[type])));
}

;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/legacy/event-modifiers.js



/**
 * Substitute for the `trusted` event modifier
 * @deprecated
 * @param {(event: Event, ...args: Array<unknown>) => void} fn
 * @returns {(event: Event, ...args: unknown[]) => void}
 */ function trusted(fn) {
    return function() {
        for(var _len = arguments.length, args = new Array(_len), _key = 0; _key < _len; _key++){
            args[_key] = arguments[_key];
        }
        var event = /** @type {Event} */ args[0];
        if (event.isTrusted) {
            // @ts-ignore
            fn === null || fn === void 0 ? void 0 : fn.apply(this, args);
        }
    };
}
/**
 * Substitute for the `self` event modifier
 * @deprecated
 * @param {(event: Event, ...args: Array<unknown>) => void} fn
 * @returns {(event: Event, ...args: unknown[]) => void}
 */ function event_modifiers_self(fn) {
    return function() {
        for(var _len = arguments.length, args = new Array(_len), _key = 0; _key < _len; _key++){
            args[_key] = arguments[_key];
        }
        var event = /** @type {Event} */ args[0];
        // @ts-ignore
        if (event.target === this) {
            // @ts-ignore
            fn === null || fn === void 0 ? void 0 : fn.apply(this, args);
        }
    };
}
/**
 * Substitute for the `stopPropagation` event modifier
 * @deprecated
 * @param {(event: Event, ...args: Array<unknown>) => void} fn
 * @returns {(event: Event, ...args: unknown[]) => void}
 */ function stopPropagation(fn) {
    return function() {
        for(var _len = arguments.length, args = new Array(_len), _key = 0; _key < _len; _key++){
            args[_key] = arguments[_key];
        }
        var event = /** @type {Event} */ args[0];
        event.stopPropagation();
        // @ts-ignore
        return fn === null || fn === void 0 ? void 0 : fn.apply(this, args);
    };
}
/**
 * Substitute for the `once` event modifier
 * @deprecated
 * @param {(event: Event, ...args: Array<unknown>) => void} fn
 * @returns {(event: Event, ...args: unknown[]) => void}
 */ function once(fn) {
    var ran = false;
    return function() {
        for(var _len = arguments.length, args = new Array(_len), _key = 0; _key < _len; _key++){
            args[_key] = arguments[_key];
        }
        if (ran) return;
        ran = true;
        // @ts-ignore
        return fn === null || fn === void 0 ? void 0 : fn.apply(this, args);
    };
}
/**
 * Substitute for the `stopImmediatePropagation` event modifier
 * @deprecated
 * @param {(event: Event, ...args: Array<unknown>) => void} fn
 * @returns {(event: Event, ...args: unknown[]) => void}
 */ function event_modifiers_stopImmediatePropagation(fn) {
    return function() {
        for(var _len = arguments.length, args = new Array(_len), _key = 0; _key < _len; _key++){
            args[_key] = arguments[_key];
        }
        var event = /** @type {Event} */ args[0];
        event.stopImmediatePropagation();
        // @ts-ignore
        return fn === null || fn === void 0 ? void 0 : fn.apply(this, args);
    };
}
/**
 * Substitute for the `preventDefault` event modifier
 * @deprecated
 * @param {(event: Event, ...args: Array<unknown>) => void} fn
 * @returns {(event: Event, ...args: unknown[]) => void}
 */ function preventDefault(fn) {
    return function() {
        for(var _len = arguments.length, args = new Array(_len), _key = 0; _key < _len; _key++){
            args[_key] = arguments[_key];
        }
        var event = /** @type {Event} */ args[0];
        event.preventDefault();
        // @ts-ignore
        return fn === null || fn === void 0 ? void 0 : fn.apply(this, args);
    };
}
/**
 * Substitute for the `passive` event modifier, implemented as an action
 * @deprecated
 * @param {HTMLElement} node
 * @param {[event: string, handler: () => EventListener]} options
 */ function passive(node, param) {
    let [event, handler] = param;
    user_pre_effect(()=>{
        return on(node, event, handler() ?? noop, {
            passive: true
        });
    });
}
/**
 * Substitute for the `nonpassive` event modifier, implemented as an action
 * @deprecated
 * @param {HTMLElement} node
 * @param {[event: string, handler: () => EventListener]} options
 */ function nonpassive(node, param) {
    let [event, handler] = param;
    user_pre_effect(()=>{
        return on(node, event, handler() ?? noop, {
            passive: false
        });
    });
}

;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/legacy/lifecycle.js
/** @import { ComponentContextLegacy } from '#client' */ 




/**
 * Legacy-mode only: Call `onMount` callbacks and set up `beforeUpdate`/`afterUpdate` effects
 * @param {boolean} [immutable]
 */ function init() {
    let immutable = arguments.length > 0 && arguments[0] !== void 0 ? arguments[0] : false;
    const context = /** @type {ComponentContextLegacy} */ component_context;
    const callbacks = context.l.u;
    if (!callbacks) return;
    let props = ()=>deep_read_state(context.s);
    if (immutable) {
        let version = 0;
        let prev = /** @type {Record<string, any>} */ {};
        // In legacy immutable mode, before/afterUpdate only fire if the object identity of a prop changes
        const d = derived(()=>{
            let changed = false;
            const props = context.s;
            for(const key in props){
                if (props[key] !== prev[key]) {
                    prev[key] = props[key];
                    changed = true;
                }
            }
            if (changed) version++;
            return version;
        });
        props = ()=>get(d);
    }
    // beforeUpdate
    if (callbacks.b.length) {
        user_pre_effect(()=>{
            observe_all(context, props);
            run_all(callbacks.b);
        });
    }
    // onMount (must run before afterUpdate)
    user_effect(()=>{
        const fns = untrack(()=>callbacks.m.map(run));
        return ()=>{
            for (const fn of fns){
                if (typeof fn === 'function') {
                    fn();
                }
            }
        };
    });
    // afterUpdate
    if (callbacks.a.length) {
        user_effect(()=>{
            observe_all(context, props);
            run_all(callbacks.a);
        });
    }
}
/**
 * Invoke the getter of all signals associated with a component
 * so they can be registered to the effect this function is called in.
 * @param {ComponentContextLegacy} context
 * @param {(() => void)} props
 */ function observe_all(context, props) {
    if (context.l.s) {
        for (const signal of context.l.s)get(signal);
    }
    props();
}

;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/legacy/misc.js



/**
 * Under some circumstances, imports may be reactive in legacy mode. In that case,
 * they should be using `reactive_import` as part of the transformation
 * @param {() => any} fn
 */ function reactive_import(fn) {
    var s = source(0);
    return function() {
        if (arguments.length === 1) {
            set(s, get(s) + 1);
            return arguments[0];
        } else {
            get(s);
            return fn();
        }
    };
}
/**
 * @this {any}
 * @param {Record<string, unknown>} $$props
 * @param {Event} event
 * @returns {void}
 */ function bubble_event($$props, event) {
    var /** @type {Record<string, Function[] | Function>} */ _$$props_$$events;
    var events = (_$$props_$$events = $$props.$$events) === null || _$$props_$$events === void 0 ? void 0 : _$$props_$$events[event.type];
    var callbacks = is_array(events) ? events.slice() : events == null ? [] : [
        events
    ];
    for (var fn of callbacks){
        // Preserve "this" context
        fn.call(this, event);
    }
}
/**
 * Used to simulate `$on` on a component instance when `compatibility.componentApi === 4`
 * @param {Record<string, any>} $$props
 * @param {string} event_name
 * @param {Function} event_callback
 */ function add_legacy_event_listener($$props, event_name, event_callback) {
    var _$$props, _$$props_$$events, _event_name;
    (_$$props = $$props).$$events || (_$$props.$$events = {});
    (_$$props_$$events = $$props.$$events)[_event_name = event_name] || (_$$props_$$events[_event_name] = []);
    $$props.$$events[event_name].push(event_callback);
}
/**
 * Used to simulate `$set` on a component instance when `compatibility.componentApi === 4`.
 * Needs component accessors so that it can call the setter of the prop. Therefore doesn't
 * work for updating props in `$$props` or `$$restProps`.
 * @this {Record<string, any>}
 * @param {Record<string, any>} $$new_props
 */ function update_legacy_props($$new_props) {
    for(var key in $$new_props){
        if (key in this) {
            this[key] = $$new_props[key];
        }
    }
}

// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/errors.js
var client_errors = __webpack_require__(626);
;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/store/utils.js
/** @import { Readable } from './public' */ 

/**
 * @template T
 * @param {Readable<T> | null | undefined} store
 * @param {(value: T) => void} run
 * @param {(value: T) => void} [invalidate]
 * @returns {() => void}
 */ function utils_subscribe_to_store(store, run, invalidate) {
    if (store == null) {
        // @ts-expect-error
        run(undefined);
        // @ts-expect-error
        if (invalidate) invalidate(undefined);
        return noop;
    }
    // Svelte store takes a private second argument
    // StartStopNotifier could mutate state, and we want to silence the corresponding validation error
    const unsub = untrack(()=>store.subscribe(run, // @ts-expect-error
        invalidate));
    // Also support RxJS
    // @ts-expect-error TODO fix this in the types?
    return unsub.unsubscribe ? ()=>unsub.unsubscribe() : unsub;
}

;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/store/shared/index.js
/** @import { Readable, StartStopNotifier, Subscriber, Unsubscriber, Updater, Writable } from '../public.js' */ /** @import { Stores, StoresValues, SubscribeInvalidateTuple } from '../private.js' */ 


/**
 * @type {Array<SubscribeInvalidateTuple<any> | any>}
 */ const subscriber_queue = (/* unused pure expression or super */ null && ([]));
/**
 * Creates a `Readable` store that allows reading by subscription.
 *
 * @template T
 * @param {T} [value] initial value
 * @param {StartStopNotifier<T>} [start]
 * @returns {Readable<T>}
 */ function readable(value, start) {
    return {
        subscribe: writable(value, start).subscribe
    };
}
/**
 * Create a `Writable` store that allows both updating and reading by subscription.
 *
 * @template T
 * @param {T} [value] initial value
 * @param {StartStopNotifier<T>} [start]
 * @returns {Writable<T>}
 */ function writable(value) {
    let start = arguments.length > 1 && arguments[1] !== void 0 ? arguments[1] : noop;
    /** @type {Unsubscriber | null} */ let stop = null;
    /** @type {Set<SubscribeInvalidateTuple<T>>} */ const subscribers = new Set();
    /**
	 * @param {T} new_value
	 * @returns {void}
	 */ function set(new_value) {
        if (safe_not_equal(value, new_value)) {
            value = new_value;
            if (stop) {
                // store is ready
                const run_queue = !subscriber_queue.length;
                for (const subscriber of subscribers){
                    subscriber[1]();
                    subscriber_queue.push(subscriber, value);
                }
                if (run_queue) {
                    for(let i = 0; i < subscriber_queue.length; i += 2){
                        subscriber_queue[i][0](subscriber_queue[i + 1]);
                    }
                    subscriber_queue.length = 0;
                }
            }
        }
    }
    /**
	 * @param {Updater<T>} fn
	 * @returns {void}
	 */ function update(fn) {
        set(fn(/** @type {T} */ value));
    }
    /**
	 * @param {Subscriber<T>} run
	 * @param {() => void} [invalidate]
	 * @returns {Unsubscriber}
	 */ function subscribe(run) {
        let invalidate = arguments.length > 1 && arguments[1] !== void 0 ? arguments[1] : noop;
        /** @type {SubscribeInvalidateTuple<T>} */ const subscriber = [
            run,
            invalidate
        ];
        subscribers.add(subscriber);
        if (subscribers.size === 1) {
            stop = start(set, update) || noop;
        }
        run(/** @type {T} */ value);
        return ()=>{
            subscribers.delete(subscriber);
            if (subscribers.size === 0 && stop) {
                stop();
                stop = null;
            }
        };
    }
    return {
        set,
        update,
        subscribe
    };
}
/**
 * Derived value store by synchronizing one or more readable stores and
 * applying an aggregation function over its input values.
 *
 * @template {Stores} S
 * @template T
 * @overload
 * @param {S} stores
 * @param {(values: StoresValues<S>, set: (value: T) => void, update: (fn: Updater<T>) => void) => Unsubscriber | void} fn
 * @param {T} [initial_value]
 * @returns {Readable<T>}
 */ /**
 * Derived value store by synchronizing one or more readable stores and
 * applying an aggregation function over its input values.
 *
 * @template {Stores} S
 * @template T
 * @overload
 * @param {S} stores
 * @param {(values: StoresValues<S>) => T} fn
 * @param {T} [initial_value]
 * @returns {Readable<T>}
 */ /**
 * @template {Stores} S
 * @template T
 * @param {S} stores
 * @param {Function} fn
 * @param {T} [initial_value]
 * @returns {Readable<T>}
 */ function shared_derived(stores, fn, initial_value) {
    const single = !Array.isArray(stores);
    /** @type {Array<Readable<any>>} */ const stores_array = single ? [
        stores
    ] : stores;
    if (!stores_array.every(Boolean)) {
        throw new Error('derived() expects stores as input, got a falsy value');
    }
    const auto = fn.length < 2;
    return readable(initial_value, (set, update)=>{
        let started = false;
        /** @type {T[]} */ const values = [];
        let pending = 0;
        let cleanup = noop;
        const sync = ()=>{
            if (pending) {
                return;
            }
            cleanup();
            const result = fn(single ? values[0] : values, set, update);
            if (auto) {
                set(result);
            } else {
                cleanup = typeof result === 'function' ? result : noop;
            }
        };
        const unsubscribers = stores_array.map((store, i)=>subscribe_to_store(store, (value)=>{
                values[i] = value;
                pending &= ~(1 << i);
                if (started) {
                    sync();
                }
            }, ()=>{
                pending |= 1 << i;
            }));
        started = true;
        sync();
        return function stop() {
            run_all(unsubscribers);
            cleanup();
            // We need to set this to false because callbacks can still happen despite having unsubscribed:
            // Callbacks might already be placed in the queue which doesn't know it should no longer
            // invoke this derived store.
            started = false;
        };
    });
}
/**
 * Takes a store and returns a new one derived from the old one that is readable.
 *
 * @template T
 * @param {Readable<T>} store  - store to make readonly
 * @returns {Readable<T>}
 */ function readonly(store) {
    return {
        // @ts-expect-error TODO i suspect the bind is unnecessary
        subscribe: store.subscribe.bind(store)
    };
}
/**
 * Get the current value from a store by subscribing and immediately unsubscribing.
 *
 * @template T
 * @param {Readable<T>} store
 * @returns {T}
 */ function shared_get(store) {
    let value;
    subscribe_to_store(store, (_)=>value = _)();
    // @ts-expect-error
    return value;
}

;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/reactivity/store.js
/** @import { StoreReferencesContainer } from '#client' */ /** @import { Store } from '#shared' */ 






/**
 * Whether or not the prop currently being read is a store binding, as in
 * `<Child bind:x={$y} />`. If it is, we treat the prop as mutable even in
 * runes mode, and skip `binding_property_non_reactive` validation
 */ let is_store_binding = false;
let IS_UNMOUNTED = Symbol();
/**
 * Gets the current value of a store. If the store isn't subscribed to yet, it will create a proxy
 * signal that will be updated when the store is. The store references container is needed to
 * track reassignments to stores and to track the correct component context.
 * @template V
 * @param {Store<V> | null | undefined} store
 * @param {string} store_name
 * @param {StoreReferencesContainer} stores
 * @returns {V}
 */ function store_get(store, store_name, stores) {
    var _stores, _store_name;
    const entry = (_stores = stores)[_store_name = store_name] ?? (_stores[_store_name] = {
        store: null,
        source: mutable_source(undefined),
        unsubscribe: noop
    });
    if (DEV) {
        entry.source.label = store_name;
    }
    // if the component that setup this is already unmounted we don't want to register a subscription
    if (entry.store !== store && !(IS_UNMOUNTED in stores)) {
        entry.unsubscribe();
        entry.store = store ?? null;
        if (store == null) {
            entry.source.v = undefined; // see synchronous callback comment below
            entry.unsubscribe = noop;
        } else {
            var is_synchronous_callback = true;
            entry.unsubscribe = subscribe_to_store(store, (v)=>{
                if (is_synchronous_callback) {
                    // If the first updates to the store value (possibly multiple of them) are synchronously
                    // inside a derived, we will hit the `state_unsafe_mutation` error if we `set` the value
                    entry.source.v = v;
                } else {
                    set(entry.source, v);
                }
            });
            is_synchronous_callback = false;
        }
    }
    // if the component that setup this stores is already unmounted the source will be out of sync
    // so we just use the `get` for the stores, less performant but it avoids to create a memory leak
    // and it will keep the value consistent
    if (store && IS_UNMOUNTED in stores) {
        return get_store(store);
    }
    return get(entry.source);
}
/**
 * Unsubscribe from a store if it's not the same as the one in the store references container.
 * We need this in addition to `store_get` because someone could unsubscribe from a store but
 * then never subscribe to the new one (if any), causing the subscription to stay open wrongfully.
 * @param {Store<any> | null | undefined} store
 * @param {string} store_name
 * @param {StoreReferencesContainer} stores
 */ function store_unsub(store, store_name, stores) {
    /** @type {StoreReferencesContainer[''] | undefined} */ let entry = stores[store_name];
    if (entry && entry.store !== store) {
        // Don't reset store yet, so that store_get above can resubscribe to new store if necessary
        entry.unsubscribe();
        entry.unsubscribe = noop;
    }
    return store;
}
/**
 * Sets the new value of a store and returns that value.
 * @template V
 * @param {Store<V>} store
 * @param {V} value
 * @returns {V}
 */ function store_set(store, value) {
    store.set(value);
    return value;
}
/**
 * @param {StoreReferencesContainer} stores
 * @param {string} store_name
 */ function invalidate_store(stores, store_name) {
    var entry = stores[store_name];
    if (entry.store !== null) {
        store_set(entry.store, entry.source.v);
    }
}
/**
 * Unsubscribes from all auto-subscribed stores on destroy
 * @returns {[StoreReferencesContainer, ()=>void]}
 */ function setup_stores() {
    /** @type {StoreReferencesContainer} */ const stores = {};
    function cleanup() {
        teardown(()=>{
            for(var store_name in stores){
                const ref = stores[store_name];
                ref.unsubscribe();
            }
            define_property(stores, IS_UNMOUNTED, {
                enumerable: false,
                value: true
            });
        });
    }
    return [
        stores,
        cleanup
    ];
}
/**
 * Updates a store with a new value.
 * @param {Store<V>} store  the store to update
 * @param {any} expression  the expression that mutates the store
 * @param {V} new_value  the new store value
 * @template V
 */ function store_mutate(store, expression, new_value) {
    store.set(new_value);
    return expression;
}
/**
 * @param {Store<number>} store
 * @param {number} store_value
 * @param {1 | -1} [d]
 * @returns {number}
 */ function update_store(store, store_value) {
    let d = arguments.length > 2 && arguments[2] !== void 0 ? arguments[2] : 1;
    store.set(store_value + d);
    return store_value;
}
/**
 * @param {Store<number>} store
 * @param {number} store_value
 * @param {1 | -1} [d]
 * @returns {number}
 */ function update_pre_store(store, store_value) {
    let d = arguments.length > 2 && arguments[2] !== void 0 ? arguments[2] : 1;
    const value = store_value + d;
    store.set(value);
    return value;
}
/**
 * Called inside prop getters to communicate that the prop is a store binding
 */ function mark_store_binding() {
    is_store_binding = true;
}
/**
 * Returns a tuple that indicates whether `fn()` reads a prop that is a store binding.
 * Used to prevent `binding_property_non_reactive` validation false positives and
 * ensure that these props are treated as mutable even in runes mode
 * @template T
 * @param {() => T} fn
 * @returns {[T, boolean]}
 */ function store_capture_store_binding(fn) {
    var previous_is_store_binding = is_store_binding;
    try {
        is_store_binding = false;
        return [
            fn(),
            is_store_binding
        ];
    } finally{
        is_store_binding = previous_is_store_binding;
    }
}

;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/reactivity/props.js
/** @import { Effect, Source } from './types.js' */ 










/**
 * @param {((value?: number) => number)} fn
 * @param {1 | -1} [d]
 * @returns {number}
 */ function update_prop(fn) {
    let d = arguments.length > 1 && arguments[1] !== void 0 ? arguments[1] : 1;
    const value = fn();
    fn(value + d);
    return value;
}
/**
 * @param {((value?: number) => number)} fn
 * @param {1 | -1} [d]
 * @returns {number}
 */ function update_pre_prop(fn) {
    let d = arguments.length > 1 && arguments[1] !== void 0 ? arguments[1] : 1;
    const value = fn() + d;
    fn(value);
    return value;
}
/**
 * The proxy handler for rest props (i.e. `const { x, ...rest } = $props()`).
 * Is passed the full `$$props` object and excludes the named props.
 * @type {ProxyHandler<{ props: Record<string | symbol, unknown>, exclude: Array<string | symbol>, name?: string }>}}
 */ const rest_props_handler = {
    get (target, key) {
        if (target.exclude.includes(key)) return;
        return target.props[key];
    },
    set (target, key) {
        if (esm_env_false/* ["default"] */.A) {
            // TODO should this happen in prod too?
            client_errors/* .props_rest_readonly */.js(`${target.name}.${String(key)}`);
        }
        return false;
    },
    getOwnPropertyDescriptor (target, key) {
        if (target.exclude.includes(key)) return;
        if (key in target.props) {
            return {
                enumerable: true,
                configurable: true,
                value: target.props[key]
            };
        }
    },
    has (target, key) {
        if (target.exclude.includes(key)) return false;
        return key in target.props;
    },
    ownKeys (target) {
        return Reflect.ownKeys(target.props).filter((key)=>!target.exclude.includes(key));
    }
};
/**
 * @param {Record<string, unknown>} props
 * @param {string[]} exclude
 * @param {string} [name]
 * @returns {Record<string, unknown>}
 */ /*#__NO_SIDE_EFFECTS__*/ function rest_props(props, exclude, name) {
    return new Proxy(esm_env_false/* ["default"] */.A ? {
        props,
        exclude,
        name,
        other: {},
        to_proxy: []
    } : {
        props,
        exclude
    }, rest_props_handler);
}
/**
 * The proxy handler for legacy $$restProps and $$props
 * @type {ProxyHandler<{ props: Record<string | symbol, unknown>, exclude: Array<string | symbol>, special: Record<string | symbol, (v?: unknown) => unknown>, version: Source<number>, parent_effect: Effect }>}}
 */ const legacy_rest_props_handler = (/* unused pure expression or super */ null && ({
    get (target, key) {
        if (target.exclude.includes(key)) return;
        get(target.version);
        return key in target.special ? target.special[key]() : target.props[key];
    },
    set (target, key, value) {
        if (!(key in target.special)) {
            var previous_effect = active_effect;
            try {
                set_active_effect(target.parent_effect);
                // Handle props that can temporarily get out of sync with the parent
                /** @type {Record<string, (v?: unknown) => unknown>} */ target.special[key] = props_prop({
                    get [key] () {
                        return target.props[key];
                    }
                }, /** @type {string} */ key, PROPS_IS_UPDATED);
            } finally{
                set_active_effect(previous_effect);
            }
        }
        target.special[key](value);
        update(target.version); // $$props is coarse-grained: when $$props.x is updated, usages of $$props.y etc are also rerun
        return true;
    },
    getOwnPropertyDescriptor (target, key) {
        if (target.exclude.includes(key)) return;
        if (key in target.props) {
            return {
                enumerable: true,
                configurable: true,
                value: target.props[key]
            };
        }
    },
    deleteProperty (target, key) {
        // Svelte 4 allowed for deletions on $$restProps
        if (target.exclude.includes(key)) return true;
        target.exclude.push(key);
        update(target.version);
        return true;
    },
    has (target, key) {
        if (target.exclude.includes(key)) return false;
        return key in target.props;
    },
    ownKeys (target) {
        return Reflect.ownKeys(target.props).filter((key)=>!target.exclude.includes(key));
    }
}));
/**
 * @param {Record<string, unknown>} props
 * @param {string[]} exclude
 * @returns {Record<string, unknown>}
 */ function legacy_rest_props(props, exclude) {
    return new Proxy({
        props,
        exclude,
        special: {},
        version: source(0),
        // TODO this is only necessary because we need to track component
        // destruction inside `prop`, because of `bind:this`, but it
        // seems likely that we can simplify `bind:this` instead
        parent_effect: /** @type {Effect} */ active_effect
    }, legacy_rest_props_handler);
}
/**
 * The proxy handler for spread props. Handles the incoming array of props
 * that looks like `() => { dynamic: props }, { static: prop }, ..` and wraps
 * them so that the whole thing is passed to the component as the `$$props` argument.
 * @type {ProxyHandler<{ props: Array<Record<string | symbol, unknown> | (() => Record<string | symbol, unknown>)> }>}}
 */ const spread_props_handler = (/* unused pure expression or super */ null && ({
    get (target, key) {
        let i = target.props.length;
        while(i--){
            let p = target.props[i];
            if (is_function(p)) p = p();
            if (typeof p === 'object' && p !== null && key in p) return p[key];
        }
    },
    set (target, key, value) {
        let i = target.props.length;
        while(i--){
            let p = target.props[i];
            if (is_function(p)) p = p();
            const desc = get_descriptor(p, key);
            if (desc && desc.set) {
                desc.set(value);
                return true;
            }
        }
        return false;
    },
    getOwnPropertyDescriptor (target, key) {
        let i = target.props.length;
        while(i--){
            let p = target.props[i];
            if (is_function(p)) p = p();
            if (typeof p === 'object' && p !== null && key in p) {
                const descriptor = get_descriptor(p, key);
                if (descriptor && !descriptor.configurable) {
                    // Prevent a "Non-configurability Report Error": The target is an array, it does
                    // not actually contain this property. If it is now described as non-configurable,
                    // the proxy throws a validation error. Setting it to true avoids that.
                    descriptor.configurable = true;
                }
                return descriptor;
            }
        }
    },
    has (target, key) {
        // To prevent a false positive `is_entry_props` in the `prop` function
        if (key === STATE_SYMBOL || key === LEGACY_PROPS) return false;
        for (let p of target.props){
            if (is_function(p)) p = p();
            if (p != null && key in p) return true;
        }
        return false;
    },
    ownKeys (target) {
        /** @type {Array<string | symbol>} */ const keys = [];
        for (let p of target.props){
            if (is_function(p)) p = p();
            if (!p) continue;
            for(const key in p){
                if (!keys.includes(key)) keys.push(key);
            }
            for (const key of Object.getOwnPropertySymbols(p)){
                if (!keys.includes(key)) keys.push(key);
            }
        }
        return keys;
    }
}));
/**
 * @param {Array<Record<string, unknown> | (() => Record<string, unknown>)>} props
 * @returns {any}
 */ function spread_props() {
    for(var _len = arguments.length, props = new Array(_len), _key = 0; _key < _len; _key++){
        props[_key] = arguments[_key];
    }
    return new Proxy({
        props
    }, spread_props_handler);
}
/**
 * This function is responsible for synchronizing a possibly bound prop with the inner component state.
 * It is used whenever the compiler sees that the component writes to the prop, or when it has a default prop_value.
 * @template V
 * @param {Record<string, unknown>} props
 * @param {string} key
 * @param {number} flags
 * @param {V | (() => V)} [fallback]
 * @returns {(() => V | ((arg: V) => V) | ((arg: V, mutation: boolean) => V))}
 */ function props_prop(props, key, flags, fallback) {
    var runes = !legacy_mode_flag || (flags & PROPS_IS_RUNES) !== 0;
    var bindable = (flags & PROPS_IS_BINDABLE) !== 0;
    var lazy = (flags & PROPS_IS_LAZY_INITIAL) !== 0;
    var fallback_value = /** @type {V} */ fallback;
    var fallback_dirty = true;
    var get_fallback = ()=>{
        if (fallback_dirty) {
            fallback_dirty = false;
            fallback_value = lazy ? untrack(/** @type {() => V} */ fallback) : /** @type {V} */ fallback;
        }
        return fallback_value;
    };
    /** @type {((v: V) => void) | undefined} */ var setter;
    if (bindable) {
        var _get_descriptor;
        // Can be the case when someone does `mount(Component, props)` with `let props = $state({...})`
        // or `createClassComponent(Component, props)`
        var is_entry_props = STATE_SYMBOL in props || LEGACY_PROPS in props;
        setter = ((_get_descriptor = get_descriptor(props, key)) === null || _get_descriptor === void 0 ? void 0 : _get_descriptor.set) ?? (is_entry_props && key in props ? (v)=>props[key] = v : undefined);
    }
    var initial_value;
    var is_store_sub = false;
    if (bindable) {
        [initial_value, is_store_sub] = capture_store_binding(()=>/** @type {V} */ props[key]);
    } else {
        initial_value = /** @type {V} */ props[key];
    }
    if (initial_value === undefined && fallback !== undefined) {
        initial_value = get_fallback();
        if (setter) {
            if (runes) e.props_invalid_value(key);
            setter(initial_value);
        }
    }
    /** @type {() => V} */ var getter;
    if (runes) {
        getter = ()=>{
            var value = /** @type {V} */ props[key];
            if (value === undefined) return get_fallback();
            fallback_dirty = true;
            return value;
        };
    } else {
        getter = ()=>{
            var value = /** @type {V} */ props[key];
            if (value !== undefined) {
                // in legacy mode, we don't revert to the fallback value
                // if the prop goes from defined to undefined. The easiest
                // way to model this is to make the fallback undefined
                // as soon as the prop has a value
                fallback_value = /** @type {V} */ undefined;
            }
            return value === undefined ? fallback_value : value;
        };
    }
    // prop is never written to — we only need a getter
    if (runes && (flags & PROPS_IS_UPDATED) === 0) {
        return getter;
    }
    // prop is written to, but the parent component had `bind:foo` which
    // means we can just call `$$props.foo = value` directly
    if (setter) {
        var legacy_parent = props.$$legacy;
        return /** @type {() => V} */ function(/** @type {V} */ value, /** @type {boolean} */ mutation) {
            if (arguments.length > 0) {
                // We don't want to notify if the value was mutated and the parent is in runes mode.
                // In that case the state proxy (if it exists) should take care of the notification.
                // If the parent is not in runes mode, we need to notify on mutation, too, that the prop
                // has changed because the parent will not be able to detect the change otherwise.
                if (!runes || !mutation || legacy_parent || is_store_sub) {
                    /** @type {Function} */ setter(mutation ? getter() : value);
                }
                return value;
            }
            return getter();
        };
    }
    // Either prop is written to, but there's no binding, which means we
    // create a derived that we can write to locally.
    // Or we are in legacy mode where we always create a derived to replicate that
    // Svelte 4 did not trigger updates when a primitive value was updated to the same value.
    var overridden = false;
    var d = ((flags & PROPS_IS_IMMUTABLE) !== 0 ? derived : derived_safe_equal)(()=>{
        overridden = false;
        return getter();
    });
    if (DEV) {
        d.label = key;
    }
    // Capture the initial value if it's bindable
    if (bindable) get(d);
    var parent_effect = /** @type {Effect} */ active_effect;
    return /** @type {() => V} */ function(/** @type {any} */ value, /** @type {boolean} */ mutation) {
        if (arguments.length > 0) {
            const new_value = mutation ? get(d) : runes && bindable ? proxy(value) : value;
            set(d, new_value);
            overridden = true;
            if (fallback_value !== undefined) {
                fallback_value = new_value;
            }
            return value;
        }
        // special case — avoid recalculating the derived if we're in a
        // teardown function and the prop was overridden locally, or the
        // component was already destroyed (this latter part is necessary
        // because `bind:this` can read props after the component has
        // been destroyed. TODO simplify `bind:this`
        if (is_destroying_effect && overridden || (parent_effect.f & DESTROYED) !== 0) {
            return d.v;
        }
        return get(d);
    };
}

// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/legacy.js
var legacy = __webpack_require__(582);
;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/validate.js







/**
 * @param {() => any} collection
 * @param {(item: any, index: number) => string} key_fn
 * @returns {void}
 */ function validate_each_keys(collection, key_fn) {
    render_effect(()=>{
        const keys = new Map();
        const maybe_array = collection();
        const array = is_array(maybe_array) ? maybe_array : maybe_array == null ? [] : Array.from(maybe_array);
        const length = array.length;
        for(let i = 0; i < length; i++){
            const key = key_fn(array[i], i);
            if (keys.has(key)) {
                const a = String(keys.get(key));
                const b = String(i);
                /** @type {string | null} */ let k = String(key);
                if (k.startsWith('[object ')) k = null;
                e.each_key_duplicate(a, b, k);
            }
            keys.set(key, i);
        }
    });
}
/**
 * @param {string} binding
 * @param {() => Record<string, any>} get_object
 * @param {() => string} get_property
 * @param {number} line
 * @param {number} column
 */ function validate_binding(binding, get_object, get_property, line, column) {
    var warned = false;
    var filename = dev_current_component_function === null || dev_current_component_function === void 0 ? void 0 : dev_current_component_function[FILENAME];
    render_effect(()=>{
        if (warned) return;
        var [object, is_store_sub] = capture_store_binding(get_object);
        if (is_store_sub) return;
        var property = get_property();
        var ran = false;
        // by making the (possibly false, but it would be an extreme edge case) assumption
        // that a getter has a corresponding setter, we can determine if a property is
        // reactive by seeing if this effect has dependencies
        var effect = render_effect(()=>{
            if (ran) return;
            // eslint-disable-next-line @typescript-eslint/no-unused-expressions
            object[property];
        });
        ran = true;
        if (effect.deps === null) {
            var location = `${filename}:${line}:${column}`;
            w.binding_property_non_reactive(binding, location);
            warned = true;
        }
    });
}

// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/flags/index.js
var internal_flags = __webpack_require__(817);
;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/legacy/legacy-client.js
/** @import { ComponentConstructorOptions, ComponentType, SvelteComponent, Component } from 'svelte' */ 















/**
 * Takes the same options as a Svelte 4 component and the component function and returns a Svelte 4 compatible component.
 *
 * @deprecated Use this only as a temporary solution to migrate your imperative component code to Svelte 5.
 *
 * @template {Record<string, any>} Props
 * @template {Record<string, any>} Exports
 * @template {Record<string, any>} Events
 * @template {Record<string, any>} Slots
 *
 * @param {ComponentConstructorOptions<Props> & {
 * 	component: ComponentType<SvelteComponent<Props, Events, Slots>> | Component<Props>;
 * }} options
 * @returns {SvelteComponent<Props, Events, Slots> & Exports}
 */ function createClassComponent(options) {
    // @ts-expect-error $$prop_def etc are not actually defined
    return new Svelte4Component(options);
}
/**
 * Takes the component function and returns a Svelte 4 compatible component constructor.
 *
 * @deprecated Use this only as a temporary solution to migrate your imperative component code to Svelte 5.
 *
 * @template {Record<string, any>} Props
 * @template {Record<string, any>} Exports
 * @template {Record<string, any>} Events
 * @template {Record<string, any>} Slots
 *
 * @param {SvelteComponent<Props, Events, Slots> | Component<Props>} component
 * @returns {ComponentType<SvelteComponent<Props, Events, Slots> & Exports>}
 */ function asClassComponent(component) {
    // @ts-expect-error $$prop_def etc are not actually defined
    return class extends Svelte4Component {
        /** @param {any} options */ constructor(options){
            super({
                component,
                ...options
            });
        }
    };
}
var /** @type {any} */ _events = /*#__PURE__*/ new WeakMap(), /** @type {Record<string, any>} */ _instance = /*#__PURE__*/ new WeakMap();
/**
 * Support using the component as both a class and function during the transition period
 * @typedef  {{new (o: ComponentConstructorOptions): SvelteComponent;(...args: Parameters<Component<Record<string, any>>>): ReturnType<Component<Record<string, any>, Record<string, any>>>;}} LegacyComponentType
 */ class Svelte4Component {
    /** @param {Record<string, any>} props */ $set(props) {
        (0,_class_private_field_get._)(this, _instance).$set(props);
    }
    /**
	 * @param {string} event
	 * @param {(...args: any[]) => any} callback
	 * @returns {any}
	 */ $on(event, callback) {
        var _this = this;
        (0,_class_private_field_get._)(this, _events)[event] = (0,_class_private_field_get._)(this, _events)[event] || [];
        /** @param {any[]} args */ const cb = function() {
            for(var _len = arguments.length, args = new Array(_len), _key = 0; _key < _len; _key++){
                args[_key] = arguments[_key];
            }
            return callback.call(_this, ...args);
        };
        (0,_class_private_field_get._)(this, _events)[event].push(cb);
        return ()=>{
            (0,_class_private_field_get._)(this, _events)[event] = (0,_class_private_field_get._)(this, _events)[event].filter(/** @param {any} fn */ (fn)=>fn !== cb);
        };
    }
    $destroy() {
        (0,_class_private_field_get._)(this, _instance).$destroy();
    }
    /**
	 * @param {ComponentConstructorOptions & {
	 *  component: any;
	 * }} options
	 */ constructor(options){
        var _options_props;
        (0,_class_private_field_init._)(this, _events, {
            writable: true,
            value: void 0
        });
        (0,_class_private_field_init._)(this, _instance, {
            writable: true,
            value: void 0
        });
        var sources = new Map();
        /**
		 * @param {string | symbol} key
		 * @param {unknown} value
		 */ var add_source = (key, value)=>{
            var s = (0,reactivity_sources/* .mutable_source */.zg)(value, false, false);
            sources.set(key, s);
            return s;
        };
        // Replicate coarse-grained props through a proxy that has a version source for
        // each property, which is incremented on updates to the property itself. Do not
        // use our $state proxy because that one has fine-grained reactivity.
        const props = new Proxy({
            ...options.props || {},
            $$events: {}
        }, {
            get (target, prop) {
                return (0,runtime/* .get */.Jt)(sources.get(prop) ?? add_source(prop, Reflect.get(target, prop)));
            },
            has (target, prop) {
                // Necessary to not throw "invalid binding" validation errors on the component side
                if (prop === client_constants/* .LEGACY_PROPS */.l3) return true;
                (0,runtime/* .get */.Jt)(sources.get(prop) ?? add_source(prop, Reflect.get(target, prop)));
                return Reflect.has(target, prop);
            },
            set (target, prop, value) {
                (0,reactivity_sources/* .set */.hZ)(sources.get(prop) ?? add_source(prop, value), value);
                return Reflect.set(target, prop, value);
            }
        });
        (0,_class_private_field_set._)(this, _instance, (options.hydrate ? render/* .hydrate */.Qv : render/* .mount */.Or)(options.component, {
            target: options.target,
            anchor: options.anchor,
            props,
            context: options.context,
            intro: options.intro ?? false,
            recover: options.recover
        }));
        // We don't flushSync for custom element wrappers or if the user doesn't want it,
        // or if we're in async mode since `flushSync()` will fail
        if (!internal_flags/* .async_mode_flag */.I0 && (!(options === null || options === void 0 ? void 0 : (_options_props = options.props) === null || _options_props === void 0 ? void 0 : _options_props.$$host) || options.sync === false)) {
            (0,reactivity_batch/* .flushSync */.qX)();
        }
        (0,_class_private_field_set._)(this, _events, props.$$events);
        for (const key of Object.keys((0,_class_private_field_get._)(this, _instance))){
            if (key === '$set' || key === '$destroy' || key === '$on') continue;
            (0,shared_utils/* .define_property */.Qu)(this, key, {
                get () {
                    return (0,_class_private_field_get._)(this, _instance)[key];
                },
                /** @param {any} value */ set (value) {
                    (0,_class_private_field_get._)(this, _instance)[key] = value;
                },
                enumerable: true
            });
        }
        (0,_class_private_field_get._)(this, _instance).$set = /** @param {Record<string, any>} next */ (next)=>{
            Object.assign(props, next);
        };
        (0,_class_private_field_get._)(this, _instance).$destroy = ()=>{
            (0,render/* .unmount */.vs)((0,_class_private_field_get._)(this, _instance));
        };
    }
}
/**
 * Runs the given function once immediately on the server, and works like `$effect.pre` on the client.
 *
 * @deprecated Use this only as a temporary solution to migrate your component code to Svelte 5.
 * @param {() => void | (() => void)} fn
 * @returns {void}
 */ function legacy_client_run(fn) {
    user_pre_effect(()=>{
        fn();
        var effect = /** @type {import('#client').Effect} */ active_effect;
        // If the effect is immediately made dirty again, mark it as maybe dirty to emulate legacy behaviour
        if ((effect.f & DIRTY) !== 0) {
            let filename = "a file (we can't know which one)";
            if (DEV) {
                // @ts-ignore
                filename = (dev_current_component_function === null || dev_current_component_function === void 0 ? void 0 : dev_current_component_function[FILENAME]) ?? filename;
            }
            w.legacy_recursive_reactive_block(filename);
            set_signal_status(effect, MAYBE_DIRTY);
        }
    });
}
/**
 * Function to mimic the multiple listeners available in svelte 4
 * @deprecated
 * @param {EventListener[]} handlers
 * @returns {EventListener}
 */ function legacy_client_handlers() {
    for(var _len = arguments.length, handlers = new Array(_len), _key = 0; _key < _len; _key++){
        handlers[_key] = arguments[_key];
    }
    return function(event) {
        const { stopImmediatePropagation } = event;
        let stopped = false;
        event.stopImmediatePropagation = ()=>{
            stopped = true;
            stopImmediatePropagation.call(event);
        };
        const errors = [];
        for (const handler of handlers){
            try {
                // @ts-expect-error `this` is not typed
                handler === null || handler === void 0 ? void 0 : handler.call(this, event);
            } catch (e) {
                errors.push(e);
            }
            if (stopped) {
                break;
            }
        }
        for (let error of errors){
            queueMicrotask(()=>{
                throw error;
            });
        }
    };
}
/**
 * Function to create a `bubble` function that mimic the behavior of `on:click` without handler available in svelte 4.
 * @deprecated Use this only as a temporary solution to migrate your automatically delegated events in Svelte 5.
 */ function createBubbler() {
    const active_component_context = component_context;
    if (active_component_context === null) {
        e.lifecycle_outside_component('createBubbler');
    }
    return (/**@type {string}*/ type)=>(/**@type {Event}*/ event)=>{
            var /** @type {Record<string, Function | Function[]>} */ _active_component_context_s_$$events;
            const events = (_active_component_context_s_$$events = active_component_context.s.$$events) === null || _active_component_context_s_$$events === void 0 ? void 0 : _active_component_context_s_$$events[/** @type {any} */ type];
            if (events) {
                const callbacks = is_array(events) ? events.slice() : [
                    events
                ];
                for (const fn of callbacks){
                    fn.call(active_component_context.x, event);
                }
                return !event.defaultPrevented;
            }
            return true;
        };
}


;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dom/elements/custom-element.js





/**
 * @typedef {Object} CustomElementPropDefinition
 * @property {string} [attribute]
 * @property {boolean} [reflect]
 * @property {'String'|'Boolean'|'Number'|'Array'|'Object'} [type]
 */ /** @type {any} */ let SvelteElement;
if (typeof HTMLElement === 'function') {
    SvelteElement = class extends HTMLElement {
        /**
		 * @param {string} type
		 * @param {EventListenerOrEventListenerObject} listener
		 * @param {boolean | AddEventListenerOptions} [options]
		 */ addEventListener(type, listener, options) {
            // We can't determine upfront if the event is a custom event or not, so we have to
            // listen to both. If someone uses a custom event with the same name as a regular
            // browser event, this fires twice - we can't avoid that.
            this.$$l[type] = this.$$l[type] || [];
            this.$$l[type].push(listener);
            if (this.$$c) {
                const unsub = this.$$c.$on(type, listener);
                this.$$l_u.set(listener, unsub);
            }
            super.addEventListener(type, listener, options);
        }
        /**
		 * @param {string} type
		 * @param {EventListenerOrEventListenerObject} listener
		 * @param {boolean | AddEventListenerOptions} [options]
		 */ removeEventListener(type, listener, options) {
            super.removeEventListener(type, listener, options);
            if (this.$$c) {
                const unsub = this.$$l_u.get(listener);
                if (unsub) {
                    unsub();
                    this.$$l_u.delete(listener);
                }
            }
        }
        async connectedCallback() {
            this.$$cn = true;
            if (!this.$$c) {
                // We wait one tick to let possible child slot elements be created/mounted
                await Promise.resolve();
                if (!this.$$cn || this.$$c) {
                    return;
                }
                /** @param {string} name */ function create_slot(name) {
                    /**
					 * @param {Element} anchor
					 */ return (anchor)=>{
                        const slot = document.createElement('slot');
                        if (name !== 'default') slot.name = name;
                        (0,template/* .append */.BC)(anchor, slot);
                    };
                }
                /** @type {Record<string, any>} */ const $$slots = {};
                const existing_slots = get_custom_elements_slots(this);
                for (const name of this.$$s){
                    if (name in existing_slots) {
                        if (name === 'default' && !this.$$d.children) {
                            this.$$d.children = create_slot(name);
                            $$slots.default = true;
                        } else {
                            $$slots[name] = create_slot(name);
                        }
                    }
                }
                for (const attribute of this.attributes){
                    // this.$$data takes precedence over this.attributes
                    const name = this.$$g_p(attribute.name);
                    if (!(name in this.$$d)) {
                        this.$$d[name] = get_custom_element_value(name, attribute.value, this.$$p_d, 'toProp');
                    }
                }
                // Port over props that were set programmatically before ce was initialized
                for(const key in this.$$p_d){
                    // @ts-expect-error
                    if (!(key in this.$$d) && this[key] !== undefined) {
                        // @ts-expect-error
                        this.$$d[key] = this[key]; // don't transform, these were set through JavaScript
                        // @ts-expect-error
                        delete this[key]; // remove the property that shadows the getter/setter
                    }
                }
                this.$$c = createClassComponent({
                    component: this.$$ctor,
                    target: this.shadowRoot || this,
                    props: {
                        ...this.$$d,
                        $$slots,
                        $$host: this
                    }
                });
                // Reflect component props as attributes
                this.$$me = (0,reactivity_effects/* .effect_root */.Fc)(()=>{
                    (0,reactivity_effects/* .render_effect */.VB)(()=>{
                        this.$$r = true;
                        for (const key of (0,shared_utils/* .object_keys */.d$)(this.$$c)){
                            var _this_$$p_d_key;
                            if (!((_this_$$p_d_key = this.$$p_d[key]) === null || _this_$$p_d_key === void 0 ? void 0 : _this_$$p_d_key.reflect)) continue;
                            this.$$d[key] = this.$$c[key];
                            const attribute_value = get_custom_element_value(key, this.$$d[key], this.$$p_d, 'toAttribute');
                            if (attribute_value == null) {
                                this.removeAttribute(this.$$p_d[key].attribute || key);
                            } else {
                                this.setAttribute(this.$$p_d[key].attribute || key, attribute_value);
                            }
                        }
                        this.$$r = false;
                    });
                });
                for(const type in this.$$l){
                    for (const listener of this.$$l[type]){
                        const unsub = this.$$c.$on(type, listener);
                        this.$$l_u.set(listener, unsub);
                    }
                }
                this.$$l = {};
            }
        }
        // We don't need this when working within Svelte code, but for compatibility of people using this outside of Svelte
        // and setting attributes through setAttribute etc, this is helpful
        /**
		 * @param {string} attr
		 * @param {string} _oldValue
		 * @param {string} newValue
		 */ attributeChangedCallback(attr, _oldValue, newValue) {
            var _this_$$c;
            if (this.$$r) return;
            attr = this.$$g_p(attr);
            this.$$d[attr] = get_custom_element_value(attr, newValue, this.$$p_d, 'toProp');
            (_this_$$c = this.$$c) === null || _this_$$c === void 0 ? void 0 : _this_$$c.$set({
                [attr]: this.$$d[attr]
            });
        }
        disconnectedCallback() {
            this.$$cn = false;
            // In a microtask, because this could be a move within the DOM
            Promise.resolve().then(()=>{
                if (!this.$$cn && this.$$c) {
                    this.$$c.$destroy();
                    this.$$me();
                    this.$$c = undefined;
                }
            });
        }
        /**
		 * @param {string} attribute_name
		 */ $$g_p(attribute_name) {
            return (0,shared_utils/* .object_keys */.d$)(this.$$p_d).find((key)=>this.$$p_d[key].attribute === attribute_name || !this.$$p_d[key].attribute && key.toLowerCase() === attribute_name) || attribute_name;
        }
        /**
		 * @param {*} $$componentCtor
		 * @param {*} $$slots
		 * @param {*} use_shadow_dom
		 */ constructor($$componentCtor, $$slots, use_shadow_dom){
            super(), /** The Svelte component constructor */ (0,_define_property._)(this, "$$ctor", void 0), /** Slots */ (0,_define_property._)(this, "$$s", void 0), /** @type {any} The Svelte component instance */ (0,_define_property._)(this, "$$c", void 0), /** Whether or not the custom element is connected */ (0,_define_property._)(this, "$$cn", false), /** @type {Record<string, any>} Component props data */ (0,_define_property._)(this, "$$d", {}), /** `true` if currently in the process of reflecting component props back to attributes */ (0,_define_property._)(this, "$$r", false), /** @type {Record<string, CustomElementPropDefinition>} Props definition (name, reflected, type etc) */ (0,_define_property._)(this, "$$p_d", {}), /** @type {Record<string, EventListenerOrEventListenerObject[]>} Event listeners */ (0,_define_property._)(this, "$$l", {}), /** @type {Map<EventListenerOrEventListenerObject, Function>} Event listener unsubscribe functions */ (0,_define_property._)(this, "$$l_u", new Map()), /** @type {any} The managed render effect for reflecting attributes */ (0,_define_property._)(this, "$$me", void 0);
            this.$$ctor = $$componentCtor;
            this.$$s = $$slots;
            if (use_shadow_dom) {
                this.attachShadow({
                    mode: 'open'
                });
            }
        }
    };
}
/**
 * @param {string} prop
 * @param {any} value
 * @param {Record<string, CustomElementPropDefinition>} props_definition
 * @param {'toAttribute' | 'toProp'} [transform]
 */ function get_custom_element_value(prop, value, props_definition, transform) {
    var _props_definition_prop;
    const type = (_props_definition_prop = props_definition[prop]) === null || _props_definition_prop === void 0 ? void 0 : _props_definition_prop.type;
    value = type === 'Boolean' && typeof value !== 'boolean' ? value != null : value;
    if (!transform || !props_definition[prop]) {
        return value;
    } else if (transform === 'toAttribute') {
        switch(type){
            case 'Object':
            case 'Array':
                return value == null ? null : JSON.stringify(value);
            case 'Boolean':
                return value ? '' : null;
            case 'Number':
                return value == null ? null : value;
            default:
                return value;
        }
    } else {
        switch(type){
            case 'Object':
            case 'Array':
                return value && JSON.parse(value);
            case 'Boolean':
                return value; // conversion already handled above
            case 'Number':
                return value != null ? +value : value;
            default:
                return value;
        }
    }
}
/**
 * @param {HTMLElement} element
 */ function get_custom_elements_slots(element) {
    /** @type {Record<string, true>} */ const result = {};
    element.childNodes.forEach((node)=>{
        result[/** @type {Element} node */ node.slot || 'default'] = true;
    });
    return result;
}
/**
 * @internal
 *
 * Turn a Svelte component into a custom element.
 * @param {any} Component  A Svelte component function
 * @param {Record<string, CustomElementPropDefinition>} props_definition  The props to observe
 * @param {string[]} slots  The slots to create
 * @param {string[]} exports  Explicitly exported values, other than props
 * @param {boolean} use_shadow_dom  Whether to use shadow DOM
 * @param {(ce: new () => HTMLElement) => new () => HTMLElement} [extend]
 */ function create_custom_element(Component, props_definition, slots, exports, use_shadow_dom, extend) {
    let Class = class extends SvelteElement {
        static get observedAttributes() {
            return object_keys(props_definition).map((key)=>(props_definition[key].attribute || key).toLowerCase());
        }
        constructor(){
            super(Component, slots, use_shadow_dom);
            this.$$p_d = props_definition;
        }
    };
    object_keys(props_definition).forEach((prop)=>{
        define_property(Class.prototype, prop, {
            get () {
                return this.$$c && prop in this.$$c ? this.$$c[prop] : this.$$d[prop];
            },
            set (value) {
                value = get_custom_element_value(prop, value, props_definition);
                this.$$d[prop] = value;
                var component = this.$$c;
                if (component) {
                    var _get_descriptor;
                    // // If the instance has an accessor, use that instead
                    var setter = (_get_descriptor = get_descriptor(component, prop)) === null || _get_descriptor === void 0 ? void 0 : _get_descriptor.get;
                    if (setter) {
                        component[prop] = value;
                    } else {
                        component.$set({
                            [prop]: value
                        });
                    }
                }
            }
        });
    });
    exports.forEach((property)=>{
        define_property(Class.prototype, property, {
            get () {
                var _this_$$c;
                return (_this_$$c = this.$$c) === null || _this_$$c === void 0 ? void 0 : _this_$$c[property];
            }
        });
    });
    if (extend) {
        // @ts-expect-error - assigning here is fine
        Class = extend(Class);
    }
    Component.element = /** @type {any} */ Class;
    return Class;
}

// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/shared/validate.js
var validate = __webpack_require__(461);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dev/equality.js
var dev_equality = __webpack_require__(301);
;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/dev/console-log.js




/**
 * @param {string} method
 * @param  {...any} objects
 */ function log_if_contains_state(method) {
    for(var _len = arguments.length, objects = new Array(_len > 1 ? _len - 1 : 0), _key = 1; _key < _len; _key++){
        objects[_key - 1] = arguments[_key];
    }
    untrack(()=>{
        try {
            let has_state = false;
            const transformed = [];
            for (const obj of objects){
                if (obj && typeof obj === 'object' && STATE_SYMBOL in obj) {
                    transformed.push(snapshot(obj, true));
                    has_state = true;
                } else {
                    transformed.push(obj);
                }
            }
            if (has_state) {
                w.console_log_state(method);
                // eslint-disable-next-line no-console
                console.log('%c[snapshot]', 'color: grey', ...transformed);
            }
        } catch  {}
    });
    return objects;
}

// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/error-handling.js
var error_handling = __webpack_require__(621);
;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/index.js









































































}),
582: (function (__unused_webpack_module, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  J: () => (captured_signals)
});
/* ESM import */var _reactivity_sources_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(264);
/* ESM import */var _runtime_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(513);
/** @import { Value } from '#client' */ 

/**
 * @type {Set<Value> | null}
 * @deprecated
 */ let captured_signals = null;
/**
 * Capture an array of all the signals that are read when `fn` is called
 * @template T
 * @param {() => T} fn
 */ function capture_signals(fn) {
    var previous_captured_signals = captured_signals;
    try {
        captured_signals = new Set();
        untrack(fn);
        if (previous_captured_signals !== null) {
            for (var signal of captured_signals){
                previous_captured_signals.add(signal);
            }
        }
        return captured_signals;
    } finally{
        captured_signals = previous_captured_signals;
    }
}
/**
 * Invokes a function and captures all signals that are read during the invocation,
 * then invalidates them.
 * @param {() => any} fn
 * @deprecated
 */ function invalidate_inner_signals(fn) {
    for (var signal of capture_signals(fn)){
        internal_set(signal, signal.v);
    }
}


}),
445: (function (__unused_webpack_module, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  B: () => (proxy),
  N: () => (get_proxied_value)
});
/* ESM import */var esm_env__WEBPACK_IMPORTED_MODULE_6__ = __webpack_require__(832);
/* ESM import */var _runtime_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(513);
/* ESM import */var _shared_utils_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(986);
/* ESM import */var _reactivity_sources_js__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(264);
/* ESM import */var _client_constants__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(924);
/* ESM import */var _constants_js__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(178);
/* ESM import */var _errors_js__WEBPACK_IMPORTED_MODULE_8__ = __webpack_require__(626);
/* ESM import */var _dev_tracing_js__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(339);
/* ESM import */var _flags_index_js__WEBPACK_IMPORTED_MODULE_7__ = __webpack_require__(817);
/** @import { Source } from '#client' */ 








// TODO move all regexes into shared module?
const regex_is_valid_identifier = /^[a-zA-Z_$][a-zA-Z_$0-9]*$/;
/**
 * @template T
 * @param {T} value
 * @returns {T}
 */ function proxy(value) {
    // if non-proxyable, or is already a proxy, return `value`
    if (typeof value !== 'object' || value === null || _client_constants__WEBPACK_IMPORTED_MODULE_3__/* .STATE_SYMBOL */.x3 in value) {
        return value;
    }
    const prototype = (0,_shared_utils_js__WEBPACK_IMPORTED_MODULE_1__/* .get_prototype_of */.Oh)(value);
    if (prototype !== _shared_utils_js__WEBPACK_IMPORTED_MODULE_1__/* .object_prototype */.N7 && prototype !== _shared_utils_js__WEBPACK_IMPORTED_MODULE_1__/* .array_prototype */.ve) {
        return value;
    }
    /** @type {Map<any, Source<any>>} */ var sources = new Map();
    var is_proxied_array = (0,_shared_utils_js__WEBPACK_IMPORTED_MODULE_1__/* .is_array */.PI)(value);
    var version = (0,_reactivity_sources_js__WEBPACK_IMPORTED_MODULE_2__/* .state */.wk)(0);
    var stack = esm_env__WEBPACK_IMPORTED_MODULE_6__/* ["default"] */.A && _flags_index_js__WEBPACK_IMPORTED_MODULE_7__/* .tracing_mode_flag */._G ? (0,_dev_tracing_js__WEBPACK_IMPORTED_MODULE_5__/* .get_stack */.sv)('CreatedAt') : null;
    var parent_version = _runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .update_version */.pJ;
    /**
	 * Executes the proxy in the context of the reaction it was originally created in, if any
	 * @template T
	 * @param {() => T} fn
	 */ var with_parent = (fn)=>{
        if (_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .update_version */.pJ === parent_version) {
            return fn();
        }
        // child source is being created after the initial proxy —
        // prevent it from being associated with the current reaction
        var reaction = _runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .active_reaction */.hp;
        var version = _runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .update_version */.pJ;
        (0,_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .set_active_reaction */.G0)(null);
        (0,_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .set_update_version */.eB)(parent_version);
        var result = fn();
        (0,_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .set_active_reaction */.G0)(reaction);
        (0,_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .set_update_version */.eB)(version);
        return result;
    };
    if (is_proxied_array) {
        // We need to create the length source eagerly to ensure that
        // mutations to the array are properly synced with our proxy
        sources.set('length', (0,_reactivity_sources_js__WEBPACK_IMPORTED_MODULE_2__/* .state */.wk)(/** @type {any[]} */ value.length, stack));
        if (esm_env__WEBPACK_IMPORTED_MODULE_6__/* ["default"] */.A) {
            value = /** @type {any} */ inspectable_array(/** @type {any[]} */ value);
        }
    }
    /** Used in dev for $inspect.trace() */ var path = '';
    let updating = false;
    /** @param {string} new_path */ function update_path(new_path) {
        if (updating) return;
        updating = true;
        path = new_path;
        (0,_dev_tracing_js__WEBPACK_IMPORTED_MODULE_5__/* .tag */.Tc)(version, `${path} version`);
        // rename all child sources and child proxies
        for (const [prop, source] of sources){
            (0,_dev_tracing_js__WEBPACK_IMPORTED_MODULE_5__/* .tag */.Tc)(source, get_label(path, prop));
        }
        updating = false;
    }
    return new Proxy(/** @type {any} */ value, {
        defineProperty (_, prop, descriptor) {
            if (!('value' in descriptor) || descriptor.configurable === false || descriptor.enumerable === false || descriptor.writable === false) {
                // we disallow non-basic descriptors, because unless they are applied to the
                // target object — which we avoid, so that state can be forked — we will run
                // afoul of the various invariants
                // https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Proxy/Proxy/getOwnPropertyDescriptor#invariants
                _errors_js__WEBPACK_IMPORTED_MODULE_8__/* .state_descriptors_fixed */.Uw();
            }
            var s = sources.get(prop);
            if (s === undefined) {
                s = with_parent(()=>{
                    var s = (0,_reactivity_sources_js__WEBPACK_IMPORTED_MODULE_2__/* .state */.wk)(descriptor.value, stack);
                    sources.set(prop, s);
                    if (esm_env__WEBPACK_IMPORTED_MODULE_6__/* ["default"] */.A && typeof prop === 'string') {
                        (0,_dev_tracing_js__WEBPACK_IMPORTED_MODULE_5__/* .tag */.Tc)(s, get_label(path, prop));
                    }
                    return s;
                });
            } else {
                (0,_reactivity_sources_js__WEBPACK_IMPORTED_MODULE_2__/* .set */.hZ)(s, descriptor.value, true);
            }
            return true;
        },
        deleteProperty (target, prop) {
            var s = sources.get(prop);
            if (s === undefined) {
                if (prop in target) {
                    const s = with_parent(()=>(0,_reactivity_sources_js__WEBPACK_IMPORTED_MODULE_2__/* .state */.wk)(_constants_js__WEBPACK_IMPORTED_MODULE_4__/* .UNINITIALIZED */.UP, stack));
                    sources.set(prop, s);
                    (0,_reactivity_sources_js__WEBPACK_IMPORTED_MODULE_2__/* .increment */.GV)(version);
                    if (esm_env__WEBPACK_IMPORTED_MODULE_6__/* ["default"] */.A) {
                        (0,_dev_tracing_js__WEBPACK_IMPORTED_MODULE_5__/* .tag */.Tc)(s, get_label(path, prop));
                    }
                }
            } else {
                (0,_reactivity_sources_js__WEBPACK_IMPORTED_MODULE_2__/* .set */.hZ)(s, _constants_js__WEBPACK_IMPORTED_MODULE_4__/* .UNINITIALIZED */.UP);
                (0,_reactivity_sources_js__WEBPACK_IMPORTED_MODULE_2__/* .increment */.GV)(version);
            }
            return true;
        },
        get (target, prop, receiver) {
            var _get_descriptor;
            if (prop === _client_constants__WEBPACK_IMPORTED_MODULE_3__/* .STATE_SYMBOL */.x3) {
                return value;
            }
            if (esm_env__WEBPACK_IMPORTED_MODULE_6__/* ["default"] */.A && prop === _client_constants__WEBPACK_IMPORTED_MODULE_3__/* .PROXY_PATH_SYMBOL */.Qf) {
                return update_path;
            }
            var s = sources.get(prop);
            var exists = prop in target;
            // create a source, but only if it's an own property and not a prototype property
            if (s === undefined && (!exists || ((_get_descriptor = (0,_shared_utils_js__WEBPACK_IMPORTED_MODULE_1__/* .get_descriptor */.J8)(target, prop)) === null || _get_descriptor === void 0 ? void 0 : _get_descriptor.writable))) {
                s = with_parent(()=>{
                    var p = proxy(exists ? target[prop] : _constants_js__WEBPACK_IMPORTED_MODULE_4__/* .UNINITIALIZED */.UP);
                    var s = (0,_reactivity_sources_js__WEBPACK_IMPORTED_MODULE_2__/* .state */.wk)(p, stack);
                    if (esm_env__WEBPACK_IMPORTED_MODULE_6__/* ["default"] */.A) {
                        (0,_dev_tracing_js__WEBPACK_IMPORTED_MODULE_5__/* .tag */.Tc)(s, get_label(path, prop));
                    }
                    return s;
                });
                sources.set(prop, s);
            }
            if (s !== undefined) {
                var v = (0,_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .get */.Jt)(s);
                return v === _constants_js__WEBPACK_IMPORTED_MODULE_4__/* .UNINITIALIZED */.UP ? undefined : v;
            }
            return Reflect.get(target, prop, receiver);
        },
        getOwnPropertyDescriptor (target, prop) {
            var descriptor = Reflect.getOwnPropertyDescriptor(target, prop);
            if (descriptor && 'value' in descriptor) {
                var s = sources.get(prop);
                if (s) descriptor.value = (0,_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .get */.Jt)(s);
            } else if (descriptor === undefined) {
                var source = sources.get(prop);
                var value = source === null || source === void 0 ? void 0 : source.v;
                if (source !== undefined && value !== _constants_js__WEBPACK_IMPORTED_MODULE_4__/* .UNINITIALIZED */.UP) {
                    return {
                        enumerable: true,
                        configurable: true,
                        value,
                        writable: true
                    };
                }
            }
            return descriptor;
        },
        has (target, prop) {
            var _get_descriptor;
            if (prop === _client_constants__WEBPACK_IMPORTED_MODULE_3__/* .STATE_SYMBOL */.x3) {
                return true;
            }
            var s = sources.get(prop);
            var has = s !== undefined && s.v !== _constants_js__WEBPACK_IMPORTED_MODULE_4__/* .UNINITIALIZED */.UP || Reflect.has(target, prop);
            if (s !== undefined || _runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .active_effect */.Fg !== null && (!has || ((_get_descriptor = (0,_shared_utils_js__WEBPACK_IMPORTED_MODULE_1__/* .get_descriptor */.J8)(target, prop)) === null || _get_descriptor === void 0 ? void 0 : _get_descriptor.writable))) {
                if (s === undefined) {
                    s = with_parent(()=>{
                        var p = has ? proxy(target[prop]) : _constants_js__WEBPACK_IMPORTED_MODULE_4__/* .UNINITIALIZED */.UP;
                        var s = (0,_reactivity_sources_js__WEBPACK_IMPORTED_MODULE_2__/* .state */.wk)(p, stack);
                        if (esm_env__WEBPACK_IMPORTED_MODULE_6__/* ["default"] */.A) {
                            (0,_dev_tracing_js__WEBPACK_IMPORTED_MODULE_5__/* .tag */.Tc)(s, get_label(path, prop));
                        }
                        return s;
                    });
                    sources.set(prop, s);
                }
                var value = (0,_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .get */.Jt)(s);
                if (value === _constants_js__WEBPACK_IMPORTED_MODULE_4__/* .UNINITIALIZED */.UP) {
                    return false;
                }
            }
            return has;
        },
        set (target, prop, value, receiver) {
            var s = sources.get(prop);
            var has = prop in target;
            // variable.length = value -> clear all signals with index >= value
            if (is_proxied_array && prop === 'length') {
                for(var i = value; i < /** @type {Source<number>} */ s.v; i += 1){
                    var other_s = sources.get(i + '');
                    if (other_s !== undefined) {
                        (0,_reactivity_sources_js__WEBPACK_IMPORTED_MODULE_2__/* .set */.hZ)(other_s, _constants_js__WEBPACK_IMPORTED_MODULE_4__/* .UNINITIALIZED */.UP);
                    } else if (i in target) {
                        // If the item exists in the original, we need to create a uninitialized source,
                        // else a later read of the property would result in a source being created with
                        // the value of the original item at that index.
                        other_s = with_parent(()=>(0,_reactivity_sources_js__WEBPACK_IMPORTED_MODULE_2__/* .state */.wk)(_constants_js__WEBPACK_IMPORTED_MODULE_4__/* .UNINITIALIZED */.UP, stack));
                        sources.set(i + '', other_s);
                        if (esm_env__WEBPACK_IMPORTED_MODULE_6__/* ["default"] */.A) {
                            (0,_dev_tracing_js__WEBPACK_IMPORTED_MODULE_5__/* .tag */.Tc)(other_s, get_label(path, i));
                        }
                    }
                }
            }
            // If we haven't yet created a source for this property, we need to ensure
            // we do so otherwise if we read it later, then the write won't be tracked and
            // the heuristics of effects will be different vs if we had read the proxied
            // object property before writing to that property.
            if (s === undefined) {
                var _get_descriptor;
                if (!has || ((_get_descriptor = (0,_shared_utils_js__WEBPACK_IMPORTED_MODULE_1__/* .get_descriptor */.J8)(target, prop)) === null || _get_descriptor === void 0 ? void 0 : _get_descriptor.writable)) {
                    s = with_parent(()=>(0,_reactivity_sources_js__WEBPACK_IMPORTED_MODULE_2__/* .state */.wk)(undefined, stack));
                    if (esm_env__WEBPACK_IMPORTED_MODULE_6__/* ["default"] */.A) {
                        (0,_dev_tracing_js__WEBPACK_IMPORTED_MODULE_5__/* .tag */.Tc)(s, get_label(path, prop));
                    }
                    (0,_reactivity_sources_js__WEBPACK_IMPORTED_MODULE_2__/* .set */.hZ)(s, proxy(value));
                    sources.set(prop, s);
                }
            } else {
                has = s.v !== _constants_js__WEBPACK_IMPORTED_MODULE_4__/* .UNINITIALIZED */.UP;
                var p = with_parent(()=>proxy(value));
                (0,_reactivity_sources_js__WEBPACK_IMPORTED_MODULE_2__/* .set */.hZ)(s, p);
            }
            var descriptor = Reflect.getOwnPropertyDescriptor(target, prop);
            // Set the new value before updating any signals so that any listeners get the new value
            if (descriptor === null || descriptor === void 0 ? void 0 : descriptor.set) {
                descriptor.set.call(receiver, value);
            }
            if (!has) {
                // If we have mutated an array directly, we might need to
                // signal that length has also changed. Do it before updating metadata
                // to ensure that iterating over the array as a result of a metadata update
                // will not cause the length to be out of sync.
                if (is_proxied_array && typeof prop === 'string') {
                    var ls = /** @type {Source<number>} */ sources.get('length');
                    var n = Number(prop);
                    if (Number.isInteger(n) && n >= ls.v) {
                        (0,_reactivity_sources_js__WEBPACK_IMPORTED_MODULE_2__/* .set */.hZ)(ls, n + 1);
                    }
                }
                (0,_reactivity_sources_js__WEBPACK_IMPORTED_MODULE_2__/* .increment */.GV)(version);
            }
            return true;
        },
        ownKeys (target) {
            (0,_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .get */.Jt)(version);
            var own_keys = Reflect.ownKeys(target).filter((key)=>{
                var source = sources.get(key);
                return source === undefined || source.v !== _constants_js__WEBPACK_IMPORTED_MODULE_4__/* .UNINITIALIZED */.UP;
            });
            for (var [key, source] of sources){
                if (source.v !== _constants_js__WEBPACK_IMPORTED_MODULE_4__/* .UNINITIALIZED */.UP && !(key in target)) {
                    own_keys.push(key);
                }
            }
            return own_keys;
        },
        setPrototypeOf () {
            _errors_js__WEBPACK_IMPORTED_MODULE_8__/* .state_prototype_fixed */.YY();
        }
    });
}
/**
 * @param {string} path
 * @param {string | symbol} prop
 */ function get_label(path, prop) {
    if (typeof prop === 'symbol') return `${path}[Symbol(${prop.description ?? ''})]`;
    if (regex_is_valid_identifier.test(prop)) return `${path}.${prop}`;
    return /^\d+$/.test(prop) ? `${path}[${prop}]` : `${path}['${prop}']`;
}
/**
 * @param {any} value
 */ function get_proxied_value(value) {
    try {
        if (value !== null && typeof value === 'object' && _client_constants__WEBPACK_IMPORTED_MODULE_3__/* .STATE_SYMBOL */.x3 in value) {
            return value[_client_constants__WEBPACK_IMPORTED_MODULE_3__/* .STATE_SYMBOL */.x3];
        }
    } catch  {
    // the above if check can throw an error if the value in question
    // is the contentWindow of an iframe on another domain, in which
    // case we want to just return the value (because it's definitely
    // not a proxied value) so we don't break any JavaScript interacting
    // with that iframe (such as various payment companies client side
    // JavaScript libraries interacting with their iframes on the same
    // domain)
    }
    return value;
}
/**
 * @param {any} a
 * @param {any} b
 */ function is(a, b) {
    return Object.is(get_proxied_value(a), get_proxied_value(b));
}
const ARRAY_MUTATING_METHODS = new Set([
    'copyWithin',
    'fill',
    'pop',
    'push',
    'reverse',
    'shift',
    'sort',
    'splice',
    'unshift'
]);
/**
 * Wrap array mutating methods so $inspect is triggered only once and
 * to prevent logging an array in intermediate state (e.g. with an empty slot)
 * @param {any[]} array
 */ function inspectable_array(array) {
    return new Proxy(array, {
        get (target, prop, receiver) {
            var value = Reflect.get(target, prop, receiver);
            if (!ARRAY_MUTATING_METHODS.has(/** @type {string} */ prop)) {
                return value;
            }
            /**
			 * @this {any[]}
			 * @param {any[]} args
			 */ return function() {
                for(var _len = arguments.length, args = new Array(_len), _key = 0; _key < _len; _key++){
                    args[_key] = arguments[_key];
                }
                (0,_reactivity_sources_js__WEBPACK_IMPORTED_MODULE_2__/* .set_inspect_effects_deferred */.HX)();
                var result = value.apply(this, args);
                (0,_reactivity_sources_js__WEBPACK_IMPORTED_MODULE_2__/* .flush_inspect_effects */.Xy)();
                return result;
            };
        }
    });
}


}),
850: (function (__unused_webpack_module, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  Bq: () => (flatten),
  sO: () => (unset_context)
});
/* ESM import */var _client_constants__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(924);
/* ESM import */var esm_env__WEBPACK_IMPORTED_MODULE_8__ = __webpack_require__(832);
/* ESM import */var _context_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(754);
/* ESM import */var _dom_blocks_boundary_js__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(899);
/* ESM import */var _error_handling_js__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(621);
/* ESM import */var _runtime_js__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(513);
/* ESM import */var _batch_js__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(410);
/* ESM import */var _deriveds_js__WEBPACK_IMPORTED_MODULE_6__ = __webpack_require__(462);
/* ESM import */var _effects_js__WEBPACK_IMPORTED_MODULE_7__ = __webpack_require__(480);
/** @import { Effect, Value } from '#client' */ 








/**
 *
 * @param {Array<() => any>} sync
 * @param {Array<() => Promise<any>>} async
 * @param {(values: Value[]) => any} fn
 */ function flatten(sync, async, fn) {
    const d = (0,_context_js__WEBPACK_IMPORTED_MODULE_1__/* .is_runes */.hH)() ? _deriveds_js__WEBPACK_IMPORTED_MODULE_6__/* .derived */.un : _deriveds_js__WEBPACK_IMPORTED_MODULE_6__/* .derived_safe_equal */.Xd;
    if (async.length === 0) {
        fn(sync.map(d));
        return;
    }
    var batch = _batch_js__WEBPACK_IMPORTED_MODULE_5__/* .current_batch */.Dr;
    var parent = /** @type {Effect} */ _runtime_js__WEBPACK_IMPORTED_MODULE_4__/* .active_effect */.Fg;
    var restore = capture();
    Promise.all(async.map((expression)=>(0,_deriveds_js__WEBPACK_IMPORTED_MODULE_6__/* .async_derived */.zx)(expression))).then((result)=>{
        batch === null || batch === void 0 ? void 0 : batch.activate();
        restore();
        try {
            fn([
                ...sync.map(d),
                ...result
            ]);
        } catch (error) {
            // ignore errors in blocks that have already been destroyed
            if ((parent.f & _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .DESTROYED */.o5) === 0) {
                (0,_error_handling_js__WEBPACK_IMPORTED_MODULE_3__/* .invoke_error_boundary */.n)(error, parent);
            }
        }
        batch === null || batch === void 0 ? void 0 : batch.deactivate();
        unset_context();
    }).catch((error)=>{
        (0,_error_handling_js__WEBPACK_IMPORTED_MODULE_3__/* .invoke_error_boundary */.n)(error, parent);
    });
}
/**
 * Captures the current effect context so that we can restore it after
 * some asynchronous work has happened (so that e.g. `await a + b`
 * causes `b` to be registered as a dependency).
 */ function capture() {
    var previous_effect = _runtime_js__WEBPACK_IMPORTED_MODULE_4__/* .active_effect */.Fg;
    var previous_reaction = _runtime_js__WEBPACK_IMPORTED_MODULE_4__/* .active_reaction */.hp;
    var previous_component_context = _context_js__WEBPACK_IMPORTED_MODULE_1__/* .component_context */.UL;
    var previous_batch = _batch_js__WEBPACK_IMPORTED_MODULE_5__/* .current_batch */.Dr;
    return function restore() {
        (0,_runtime_js__WEBPACK_IMPORTED_MODULE_4__/* .set_active_effect */.gU)(previous_effect);
        (0,_runtime_js__WEBPACK_IMPORTED_MODULE_4__/* .set_active_reaction */.G0)(previous_reaction);
        (0,_context_js__WEBPACK_IMPORTED_MODULE_1__/* .set_component_context */.De)(previous_component_context);
        previous_batch === null || previous_batch === void 0 ? void 0 : previous_batch.activate();
        if (esm_env__WEBPACK_IMPORTED_MODULE_8__/* ["default"] */.A) {
            (0,_deriveds_js__WEBPACK_IMPORTED_MODULE_6__/* .set_from_async_derived */.V)(null);
        }
    };
}
/**
 * Wraps an `await` expression in such a way that the effect context that was
 * active before the expression evaluated can be reapplied afterwards —
 * `await a + b` becomes `(await $.save(a))() + b`
 * @template T
 * @param {Promise<T>} promise
 * @returns {Promise<() => T>}
 */ async function save(promise) {
    var restore = capture();
    var value = await promise;
    return ()=>{
        restore();
        return value;
    };
}
/**
 * Reset `current_async_effect` after the `promise` resolves, so
 * that we can emit `await_reactivity_loss` warnings
 * @template T
 * @param {Promise<T>} promise
 * @returns {Promise<() => T>}
 */ async function track_reactivity_loss(promise) {
    var previous_async_effect = current_async_effect;
    var value = await promise;
    return ()=>{
        set_from_async_derived(previous_async_effect);
        return value;
    };
}
/**
 * Used in `for await` loops in DEV, so
 * that we can emit `await_reactivity_loss` warnings
 * after each `async_iterator` result resolves and
 * after the `async_iterator` return resolves (if it runs)
 * @template T
 * @template TReturn
 * @param {Iterable<T> | AsyncIterable<T>} iterable
 * @returns {AsyncGenerator<T, TReturn | undefined>}
 */ async function* for_await_track_reactivity_loss(iterable) {
    var _iterable_Symbol_asyncIterator, _iterable_Symbol_iterator;
    // This is based on the algorithms described in ECMA-262:
    // ForIn/OfBodyEvaluation
    // https://tc39.es/ecma262/multipage/ecmascript-language-statements-and-declarations.html#sec-runtime-semantics-forin-div-ofbodyevaluation-lhs-stmt-iterator-lhskind-labelset
    // AsyncIteratorClose
    // https://tc39.es/ecma262/multipage/abstract-operations.html#sec-asynciteratorclose
    /** @type {AsyncIterator<T, TReturn>} */ // @ts-ignore
    const iterator = ((_iterable_Symbol_asyncIterator = iterable[Symbol.asyncIterator]) === null || _iterable_Symbol_asyncIterator === void 0 ? void 0 : _iterable_Symbol_asyncIterator.call(iterable)) ?? ((_iterable_Symbol_iterator = iterable[Symbol.iterator]) === null || _iterable_Symbol_iterator === void 0 ? void 0 : _iterable_Symbol_iterator.call(iterable));
    if (iterator === undefined) {
        throw new TypeError('value is not async iterable');
    }
    /** Whether the completion of the iterator was "normal", meaning it wasn't ended via `break` or a similar method */ let normal_completion = false;
    try {
        while(true){
            const { done, value } = (await track_reactivity_loss(iterator.next()))();
            if (done) {
                normal_completion = true;
                break;
            }
            yield value;
        }
    } finally{
        // If the iterator had a normal completion and `return` is defined on the iterator, call it and return the value
        if (normal_completion && iterator.return !== undefined) {
            // eslint-disable-next-line no-unsafe-finally
            return /** @type {TReturn} */ (await track_reactivity_loss(iterator.return()))().value;
        }
    }
}
function unset_context() {
    (0,_runtime_js__WEBPACK_IMPORTED_MODULE_4__/* .set_active_effect */.gU)(null);
    (0,_runtime_js__WEBPACK_IMPORTED_MODULE_4__/* .set_active_reaction */.G0)(null);
    (0,_context_js__WEBPACK_IMPORTED_MODULE_1__/* .set_component_context */.De)(null);
    if (esm_env__WEBPACK_IMPORTED_MODULE_8__/* ["default"] */.A) (0,_deriveds_js__WEBPACK_IMPORTED_MODULE_6__/* .set_from_async_derived */.V)(null);
}
/**
 * @param {() => Promise<void>} fn
 */ async function async_body(fn) {
    var boundary = get_boundary();
    var batch = /** @type {Batch} */ current_batch;
    var pending = boundary.is_pending();
    boundary.update_pending_count(1);
    if (!pending) batch.increment();
    var active = /** @type {Effect} */ active_effect;
    try {
        await fn();
    } catch (error) {
        if (!aborted(active)) {
            invoke_error_boundary(error, active);
        }
    } finally{
        boundary.update_pending_count(-1);
        if (pending) {
            batch.flush();
        } else {
            batch.activate();
            batch.decrement();
        }
        unset_context();
    }
}


}),
410: (function (__unused_webpack_module, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  Dr: () => (current_batch),
  G1: () => (batch_deriveds),
  OH: () => (is_flushing_sync),
  ec: () => (schedule_effect),
  es: () => (eager_block_effects),
  lP: () => (Batch),
  qX: () => (flushSync),
  x8: () => (effect_pending_updates)
});
/* ESM import */var _swc_helpers_class_private_field_get__WEBPACK_IMPORTED_MODULE_8__ = __webpack_require__(570);
/* ESM import */var _swc_helpers_class_private_field_init__WEBPACK_IMPORTED_MODULE_13__ = __webpack_require__(636);
/* ESM import */var _swc_helpers_class_private_field_set__WEBPACK_IMPORTED_MODULE_10__ = __webpack_require__(549);
/* ESM import */var _swc_helpers_class_private_method_get__WEBPACK_IMPORTED_MODULE_9__ = __webpack_require__(585);
/* ESM import */var _swc_helpers_class_private_method_init__WEBPACK_IMPORTED_MODULE_11__ = __webpack_require__(23);
/* ESM import */var _swc_helpers_define_property__WEBPACK_IMPORTED_MODULE_12__ = __webpack_require__(925);
/* ESM import */var _client_constants__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(924);
/* ESM import */var _flags_index_js__WEBPACK_IMPORTED_MODULE_7__ = __webpack_require__(817);
/* ESM import */var _shared_utils_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(986);
/* ESM import */var _runtime_js__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(513);
/* ESM import */var _errors_js__WEBPACK_IMPORTED_MODULE_14__ = __webpack_require__(626);
/* ESM import */var _dom_task_js__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(593);
/* ESM import */var esm_env__WEBPACK_IMPORTED_MODULE_15__ = __webpack_require__(832);
/* ESM import */var _error_handling_js__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(621);
/* ESM import */var _sources_js__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(264);
/* ESM import */var _effects_js__WEBPACK_IMPORTED_MODULE_6__ = __webpack_require__(480);
/** @import { Derived, Effect, Source } from '#client' */ 















/** @type {Set<Batch>} */ const batches = new Set();
/** @type {Batch | null} */ let current_batch = null;
/**
 * This is needed to avoid overwriting inputs in non-async mode
 * TODO 6.0 remove this, as non-async mode will go away
 * @type {Batch | null}
 */ let previous_batch = null;
/**
 * When time travelling, we re-evaluate deriveds based on the temporary
 * values of their dependencies rather than their actual values, and cache
 * the results in this map rather than on the deriveds themselves
 * @type {Map<Derived, any> | null}
 */ let batch_deriveds = null;
/** @type {Set<() => void>} */ let effect_pending_updates = new Set();
/** @type {Effect[]} */ let queued_root_effects = [];
/** @type {Effect | null} */ let last_scheduled_effect = null;
let is_flushing = false;
let is_flushing_sync = false;
var /**
	 * The values of any sources that are updated in this batch _before_ those updates took place.
	 * They keys of this map are identical to `this.#current`
	 * @type {Map<Source, any>}
	 */ _previous = /*#__PURE__*/ new WeakMap(), /**
	 * When the batch is committed (and the DOM is updated), we need to remove old branches
	 * and append new ones by calling the functions added inside (if/each/key/etc) blocks
	 * @type {Set<() => void>}
	 */ _callbacks = /*#__PURE__*/ new WeakMap(), /**
	 * The number of async effects that are currently in flight
	 */ _pending = /*#__PURE__*/ new WeakMap(), /**
	 * A deferred that resolves when the batch is committed, used with `settled()`
	 * TODO replace with Promise.withResolvers once supported widely enough
	 * @type {{ promise: Promise<void>, resolve: (value?: any) => void, reject: (reason: unknown) => void } | null}
	 */ _deferred = /*#__PURE__*/ new WeakMap(), /**
	 * True if an async effect inside this batch resolved and
	 * its parent branch was already deleted
	 */ _neutered = /*#__PURE__*/ new WeakMap(), /**
	 * Async effects (created inside `async_derived`) encountered during processing.
	 * These run after the rest of the batch has updated, since they should
	 * always have the latest values
	 * @type {Effect[]}
	 */ _async_effects = /*#__PURE__*/ new WeakMap(), /**
	 * The same as `#async_effects`, but for effects inside a newly-created
	 * `<svelte:boundary>` — these do not prevent the batch from committing
	 * @type {Effect[]}
	 */ _boundary_async_effects = /*#__PURE__*/ new WeakMap(), /**
	 * Template effects and `$effect.pre` effects, which run when
	 * a batch is committed
	 * @type {Effect[]}
	 */ _render_effects = /*#__PURE__*/ new WeakMap(), /**
	 * The same as `#render_effects`, but for `$effect` (which runs after)
	 * @type {Effect[]}
	 */ _effects = /*#__PURE__*/ new WeakMap(), /**
	 * Block effects, which may need to re-run on subsequent flushes
	 * in order to update internal sources (e.g. each block items)
	 * @type {Effect[]}
	 */ _block_effects = /*#__PURE__*/ new WeakMap(), /**
	 * Deferred effects (which run after async work has completed) that are DIRTY
	 * @type {Effect[]}
	 */ _dirty_effects = /*#__PURE__*/ new WeakMap(), /**
	 * Deferred effects that are MAYBE_DIRTY
	 * @type {Effect[]}
	 */ _maybe_dirty_effects = /*#__PURE__*/ new WeakMap(), /**
	 * Traverse the effect tree, executing effects or stashing
	 * them for later execution as appropriate
	 * @param {Effect} root
	 */ _traverse_effect_tree = /*#__PURE__*/ new WeakSet(), /**
	 * @param {Effect[]} effects
	 */ _defer_effects = /*#__PURE__*/ new WeakSet(), /**
	 * Append and remove branches to/from the DOM
	 */ _commit = /*#__PURE__*/ new WeakSet();
class Batch {
    /**
	 *
	 * @param {Effect[]} root_effects
	 */ process(root_effects) {
        queued_root_effects = [];
        previous_batch = null;
        /** @type {Map<Source, { v: unknown, wv: number }> | null} */ var current_values = null;
        // if there are multiple batches, we are 'time travelling' —
        // we need to undo the changes belonging to any batch
        // other than the current one
        if (_flags_index_js__WEBPACK_IMPORTED_MODULE_7__/* .async_mode_flag */.I0 && batches.size > 1) {
            current_values = new Map();
            batch_deriveds = new Map();
            for (const [source, current] of this.current){
                current_values.set(source, {
                    v: source.v,
                    wv: source.wv
                });
                source.v = current;
            }
            for (const batch of batches){
                if (batch === this) continue;
                for (const [source, previous] of (0,_swc_helpers_class_private_field_get__WEBPACK_IMPORTED_MODULE_8__._)(batch, _previous)){
                    if (!current_values.has(source)) {
                        current_values.set(source, {
                            v: source.v,
                            wv: source.wv
                        });
                        source.v = previous;
                    }
                }
            }
        }
        for (const root of root_effects){
            (0,_swc_helpers_class_private_method_get__WEBPACK_IMPORTED_MODULE_9__._)(this, _traverse_effect_tree, traverse_effect_tree).call(this, root);
        }
        // if we didn't start any new async work, and no async work
        // is outstanding from a previous flush, commit
        if ((0,_swc_helpers_class_private_field_get__WEBPACK_IMPORTED_MODULE_8__._)(this, _async_effects).length === 0 && (0,_swc_helpers_class_private_field_get__WEBPACK_IMPORTED_MODULE_8__._)(this, _pending) === 0) {
            var _class_private_field_get1;
            (0,_swc_helpers_class_private_method_get__WEBPACK_IMPORTED_MODULE_9__._)(this, _commit, commit).call(this);
            var render_effects = (0,_swc_helpers_class_private_field_get__WEBPACK_IMPORTED_MODULE_8__._)(this, _render_effects);
            var effects = (0,_swc_helpers_class_private_field_get__WEBPACK_IMPORTED_MODULE_8__._)(this, _effects);
            (0,_swc_helpers_class_private_field_set__WEBPACK_IMPORTED_MODULE_10__._)(this, _render_effects, []);
            (0,_swc_helpers_class_private_field_set__WEBPACK_IMPORTED_MODULE_10__._)(this, _effects, []);
            (0,_swc_helpers_class_private_field_set__WEBPACK_IMPORTED_MODULE_10__._)(this, _block_effects, []);
            // If sources are written to, then work needs to happen in a separate batch, else prior sources would be mixed with
            // newly updated sources, which could lead to infinite loops when effects run over and over again.
            previous_batch = current_batch;
            current_batch = null;
            flush_queued_effects(render_effects);
            flush_queued_effects(effects);
            // Reinstate the current batch if there was no new one created, as `process()` runs in a loop in `flush_effects()`.
            // That method expects `current_batch` to be set, and could run the loop again if effects result in new effects
            // being scheduled but without writes happening in which case no new batch is created.
            if (current_batch === null) {
                current_batch = this;
            } else {
                batches.delete(this);
            }
            (_class_private_field_get1 = (0,_swc_helpers_class_private_field_get__WEBPACK_IMPORTED_MODULE_8__._)(this, _deferred)) === null || _class_private_field_get1 === void 0 ? void 0 : _class_private_field_get1.resolve();
        } else {
            (0,_swc_helpers_class_private_method_get__WEBPACK_IMPORTED_MODULE_9__._)(this, _defer_effects, defer_effects).call(this, (0,_swc_helpers_class_private_field_get__WEBPACK_IMPORTED_MODULE_8__._)(this, _render_effects));
            (0,_swc_helpers_class_private_method_get__WEBPACK_IMPORTED_MODULE_9__._)(this, _defer_effects, defer_effects).call(this, (0,_swc_helpers_class_private_field_get__WEBPACK_IMPORTED_MODULE_8__._)(this, _effects));
            (0,_swc_helpers_class_private_method_get__WEBPACK_IMPORTED_MODULE_9__._)(this, _defer_effects, defer_effects).call(this, (0,_swc_helpers_class_private_field_get__WEBPACK_IMPORTED_MODULE_8__._)(this, _block_effects));
        }
        if (current_values) {
            for (const [source, { v, wv }] of current_values){
                // reset the source to the current value (unless
                // it got a newer value as a result of effects running)
                if (source.wv <= wv) {
                    source.v = v;
                }
            }
            batch_deriveds = null;
        }
        for (const effect of (0,_swc_helpers_class_private_field_get__WEBPACK_IMPORTED_MODULE_8__._)(this, _async_effects)){
            (0,_runtime_js__WEBPACK_IMPORTED_MODULE_2__/* .update_effect */.gJ)(effect);
        }
        for (const effect of (0,_swc_helpers_class_private_field_get__WEBPACK_IMPORTED_MODULE_8__._)(this, _boundary_async_effects)){
            (0,_runtime_js__WEBPACK_IMPORTED_MODULE_2__/* .update_effect */.gJ)(effect);
        }
        (0,_swc_helpers_class_private_field_set__WEBPACK_IMPORTED_MODULE_10__._)(this, _async_effects, []);
        (0,_swc_helpers_class_private_field_set__WEBPACK_IMPORTED_MODULE_10__._)(this, _boundary_async_effects, []);
    }
    /**
	 * Associate a change to a given source with the current
	 * batch, noting its previous and current values
	 * @param {Source} source
	 * @param {any} value
	 */ capture(source, value) {
        if (!(0,_swc_helpers_class_private_field_get__WEBPACK_IMPORTED_MODULE_8__._)(this, _previous).has(source)) {
            (0,_swc_helpers_class_private_field_get__WEBPACK_IMPORTED_MODULE_8__._)(this, _previous).set(source, value);
        }
        this.current.set(source, source.v);
    }
    activate() {
        current_batch = this;
    }
    deactivate() {
        current_batch = null;
        previous_batch = null;
        for (const update of effect_pending_updates){
            effect_pending_updates.delete(update);
            update();
            if (current_batch !== null) {
                break;
            }
        }
    }
    neuter() {
        (0,_swc_helpers_class_private_field_set__WEBPACK_IMPORTED_MODULE_10__._)(this, _neutered, true);
    }
    flush() {
        if (queued_root_effects.length > 0) {
            flush_effects();
        } else {
            (0,_swc_helpers_class_private_method_get__WEBPACK_IMPORTED_MODULE_9__._)(this, _commit, commit).call(this);
        }
        if (current_batch !== this) {
            // this can happen if a `flushSync` occurred during `flush_effects()`,
            // which is permitted in legacy mode despite being a terrible idea
            return;
        }
        if ((0,_swc_helpers_class_private_field_get__WEBPACK_IMPORTED_MODULE_8__._)(this, _pending) === 0) {
            batches.delete(this);
        }
        this.deactivate();
    }
    increment() {
        (0,_swc_helpers_class_private_field_set__WEBPACK_IMPORTED_MODULE_10__._)(this, _pending, (0,_swc_helpers_class_private_field_get__WEBPACK_IMPORTED_MODULE_8__._)(this, _pending) + 1);
    }
    decrement() {
        (0,_swc_helpers_class_private_field_set__WEBPACK_IMPORTED_MODULE_10__._)(this, _pending, (0,_swc_helpers_class_private_field_get__WEBPACK_IMPORTED_MODULE_8__._)(this, _pending) - 1);
        if ((0,_swc_helpers_class_private_field_get__WEBPACK_IMPORTED_MODULE_8__._)(this, _pending) === 0) {
            for (const e of (0,_swc_helpers_class_private_field_get__WEBPACK_IMPORTED_MODULE_8__._)(this, _dirty_effects)){
                (0,_runtime_js__WEBPACK_IMPORTED_MODULE_2__/* .set_signal_status */.TC)(e, _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .DIRTY */.jm);
                schedule_effect(e);
            }
            for (const e of (0,_swc_helpers_class_private_field_get__WEBPACK_IMPORTED_MODULE_8__._)(this, _maybe_dirty_effects)){
                (0,_runtime_js__WEBPACK_IMPORTED_MODULE_2__/* .set_signal_status */.TC)(e, _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .MAYBE_DIRTY */.ig);
                schedule_effect(e);
            }
            (0,_swc_helpers_class_private_field_set__WEBPACK_IMPORTED_MODULE_10__._)(this, _render_effects, []);
            (0,_swc_helpers_class_private_field_set__WEBPACK_IMPORTED_MODULE_10__._)(this, _effects, []);
            this.flush();
        } else {
            this.deactivate();
        }
    }
    /** @param {() => void} fn */ add_callback(fn) {
        (0,_swc_helpers_class_private_field_get__WEBPACK_IMPORTED_MODULE_8__._)(this, _callbacks).add(fn);
    }
    settled() {
        return (0,_swc_helpers_class_private_field_set__WEBPACK_IMPORTED_MODULE_10__._)(this, _deferred, (0,_swc_helpers_class_private_field_get__WEBPACK_IMPORTED_MODULE_8__._)(this, _deferred) ?? (0,_shared_utils_js__WEBPACK_IMPORTED_MODULE_1__/* .deferred */.yX)()).promise;
    }
    static ensure() {
        if (current_batch === null) {
            const batch = current_batch = new Batch();
            batches.add(current_batch);
            if (!is_flushing_sync) {
                Batch.enqueue(()=>{
                    if (current_batch !== batch) {
                        // a flushSync happened in the meantime
                        return;
                    }
                    batch.flush();
                });
            }
        }
        return current_batch;
    }
    /** @param {() => void} task */ static enqueue(task) {
        (0,_dom_task_js__WEBPACK_IMPORTED_MODULE_3__/* .queue_micro_task */.$r)(task);
    }
    constructor(){
        (0,_swc_helpers_class_private_method_init__WEBPACK_IMPORTED_MODULE_11__._)(this, _traverse_effect_tree);
        (0,_swc_helpers_class_private_method_init__WEBPACK_IMPORTED_MODULE_11__._)(this, _defer_effects);
        (0,_swc_helpers_class_private_method_init__WEBPACK_IMPORTED_MODULE_11__._)(this, _commit);
        /**
	 * The current values of any sources that are updated in this batch
	 * They keys of this map are identical to `this.#previous`
	 * @type {Map<Source, any>}
	 */ (0,_swc_helpers_define_property__WEBPACK_IMPORTED_MODULE_12__._)(this, "current", new Map());
        (0,_swc_helpers_class_private_field_init__WEBPACK_IMPORTED_MODULE_13__._)(this, _previous, {
            writable: true,
            value: new Map()
        });
        (0,_swc_helpers_class_private_field_init__WEBPACK_IMPORTED_MODULE_13__._)(this, _callbacks, {
            writable: true,
            value: new Set()
        });
        (0,_swc_helpers_class_private_field_init__WEBPACK_IMPORTED_MODULE_13__._)(this, _pending, {
            writable: true,
            value: 0
        });
        (0,_swc_helpers_class_private_field_init__WEBPACK_IMPORTED_MODULE_13__._)(this, _deferred, {
            writable: true,
            value: null
        });
        (0,_swc_helpers_class_private_field_init__WEBPACK_IMPORTED_MODULE_13__._)(this, _neutered, {
            writable: true,
            value: false
        });
        (0,_swc_helpers_class_private_field_init__WEBPACK_IMPORTED_MODULE_13__._)(this, _async_effects, {
            writable: true,
            value: []
        });
        (0,_swc_helpers_class_private_field_init__WEBPACK_IMPORTED_MODULE_13__._)(this, _boundary_async_effects, {
            writable: true,
            value: []
        });
        (0,_swc_helpers_class_private_field_init__WEBPACK_IMPORTED_MODULE_13__._)(this, _render_effects, {
            writable: true,
            value: []
        });
        (0,_swc_helpers_class_private_field_init__WEBPACK_IMPORTED_MODULE_13__._)(this, _effects, {
            writable: true,
            value: []
        });
        (0,_swc_helpers_class_private_field_init__WEBPACK_IMPORTED_MODULE_13__._)(this, _block_effects, {
            writable: true,
            value: []
        });
        (0,_swc_helpers_class_private_field_init__WEBPACK_IMPORTED_MODULE_13__._)(this, _dirty_effects, {
            writable: true,
            value: []
        });
        (0,_swc_helpers_class_private_field_init__WEBPACK_IMPORTED_MODULE_13__._)(this, _maybe_dirty_effects, {
            writable: true,
            value: []
        });
        /**
	 * A set of branches that still exist, but will be destroyed when this batch
	 * is committed — we skip over these during `process`
	 * @type {Set<Effect>}
	 */ (0,_swc_helpers_define_property__WEBPACK_IMPORTED_MODULE_12__._)(this, "skipped_effects", new Set());
    }
}
function traverse_effect_tree(root) {
    root.f ^= _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .CLEAN */.w_;
    var effect = root.first;
    while(effect !== null){
        var flags = effect.f;
        var is_branch = (flags & (_client_constants__WEBPACK_IMPORTED_MODULE_0__/* .BRANCH_EFFECT */.Zr | _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .ROOT_EFFECT */.FV)) !== 0;
        var is_skippable_branch = is_branch && (flags & _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .CLEAN */.w_) !== 0;
        var skip = is_skippable_branch || (flags & _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .INERT */.$q) !== 0 || this.skipped_effects.has(effect);
        if (!skip && effect.fn !== null) {
            if (is_branch) {
                effect.f ^= _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .CLEAN */.w_;
            } else if ((flags & _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .EFFECT */.ac) !== 0) {
                (0,_swc_helpers_class_private_field_get__WEBPACK_IMPORTED_MODULE_8__._)(this, _effects).push(effect);
            } else if (_flags_index_js__WEBPACK_IMPORTED_MODULE_7__/* .async_mode_flag */.I0 && (flags & _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .RENDER_EFFECT */.Zv) !== 0) {
                (0,_swc_helpers_class_private_field_get__WEBPACK_IMPORTED_MODULE_8__._)(this, _render_effects).push(effect);
            } else if ((flags & _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .CLEAN */.w_) === 0) {
                if ((flags & _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .ASYNC */.VD) !== 0) {
                    var _effect_b;
                    var effects = ((_effect_b = effect.b) === null || _effect_b === void 0 ? void 0 : _effect_b.is_pending()) ? (0,_swc_helpers_class_private_field_get__WEBPACK_IMPORTED_MODULE_8__._)(this, _boundary_async_effects) : (0,_swc_helpers_class_private_field_get__WEBPACK_IMPORTED_MODULE_8__._)(this, _async_effects);
                    effects.push(effect);
                } else if ((0,_runtime_js__WEBPACK_IMPORTED_MODULE_2__/* .is_dirty */.Kj)(effect)) {
                    if ((effect.f & _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .BLOCK_EFFECT */.kc) !== 0) (0,_swc_helpers_class_private_field_get__WEBPACK_IMPORTED_MODULE_8__._)(this, _block_effects).push(effect);
                    (0,_runtime_js__WEBPACK_IMPORTED_MODULE_2__/* .update_effect */.gJ)(effect);
                }
            }
            var child = effect.first;
            if (child !== null) {
                effect = child;
                continue;
            }
        }
        var parent = effect.parent;
        effect = effect.next;
        while(effect === null && parent !== null){
            effect = parent.next;
            parent = parent.parent;
        }
    }
}
function defer_effects(effects) {
    for (const e of effects){
        const target = (e.f & _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .DIRTY */.jm) !== 0 ? (0,_swc_helpers_class_private_field_get__WEBPACK_IMPORTED_MODULE_8__._)(this, _dirty_effects) : (0,_swc_helpers_class_private_field_get__WEBPACK_IMPORTED_MODULE_8__._)(this, _maybe_dirty_effects);
        target.push(e);
        // mark as clean so they get scheduled if they depend on pending async state
        (0,_runtime_js__WEBPACK_IMPORTED_MODULE_2__/* .set_signal_status */.TC)(e, _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .CLEAN */.w_);
    }
    effects.length = 0;
}
function commit() {
    if (!(0,_swc_helpers_class_private_field_get__WEBPACK_IMPORTED_MODULE_8__._)(this, _neutered)) {
        for (const fn of (0,_swc_helpers_class_private_field_get__WEBPACK_IMPORTED_MODULE_8__._)(this, _callbacks)){
            fn();
        }
    }
    (0,_swc_helpers_class_private_field_get__WEBPACK_IMPORTED_MODULE_8__._)(this, _callbacks).clear();
}
/**
 * Synchronously flush any pending updates.
 * Returns void if no callback is provided, otherwise returns the result of calling the callback.
 * @template [T=void]
 * @param {(() => T) | undefined} [fn]
 * @returns {T}
 */ function flushSync(fn) {
    if (_flags_index_js__WEBPACK_IMPORTED_MODULE_7__/* .async_mode_flag */.I0 && _runtime_js__WEBPACK_IMPORTED_MODULE_2__/* .active_effect */.Fg !== null) {
        // We disallow this because it creates super-hard to reason about stack trace and because it's generally a bad idea
        _errors_js__WEBPACK_IMPORTED_MODULE_14__/* .flush_sync_in_effect */.fW();
    }
    var was_flushing_sync = is_flushing_sync;
    is_flushing_sync = true;
    try {
        var result;
        if (fn) {
            flush_effects();
            result = fn();
        }
        while(true){
            (0,_dom_task_js__WEBPACK_IMPORTED_MODULE_3__/* .flush_tasks */.eo)();
            if (queued_root_effects.length === 0 && !(0,_dom_task_js__WEBPACK_IMPORTED_MODULE_3__/* .has_pending_tasks */.KJ)()) {
                current_batch === null || current_batch === void 0 ? void 0 : current_batch.flush();
                // we need to check again, in case we just updated an `$effect.pending()`
                if (queued_root_effects.length === 0) {
                    // this would be reset in `flush_effects()` but since we are early returning here,
                    // we need to reset it here as well in case the first time there's 0 queued root effects
                    last_scheduled_effect = null;
                    return /** @type {T} */ result;
                }
            }
            flush_effects();
        }
    } finally{
        is_flushing_sync = was_flushing_sync;
    }
}
function flush_effects() {
    var was_updating_effect = _runtime_js__WEBPACK_IMPORTED_MODULE_2__/* .is_updating_effect */.st;
    is_flushing = true;
    try {
        var flush_count = 0;
        (0,_runtime_js__WEBPACK_IMPORTED_MODULE_2__/* .set_is_updating_effect */.BI)(true);
        while(queued_root_effects.length > 0){
            var batch = Batch.ensure();
            if (flush_count++ > 1000) {
                if (esm_env__WEBPACK_IMPORTED_MODULE_15__/* ["default"] */.A) {
                    var updates = new Map();
                    for (const source of batch.current.keys()){
                        for (const [stack, update] of source.updated ?? []){
                            var entry = updates.get(stack);
                            if (!entry) {
                                entry = {
                                    error: update.error,
                                    count: 0
                                };
                                updates.set(stack, entry);
                            }
                            entry.count += update.count;
                        }
                    }
                    for (const update of updates.values()){
                        // eslint-disable-next-line no-console
                        console.error(update.error);
                    }
                }
                infinite_loop_guard();
            }
            batch.process(queued_root_effects);
            _sources_js__WEBPACK_IMPORTED_MODULE_5__/* .old_values.clear */.bJ.clear();
        }
    } finally{
        is_flushing = false;
        (0,_runtime_js__WEBPACK_IMPORTED_MODULE_2__/* .set_is_updating_effect */.BI)(was_updating_effect);
        last_scheduled_effect = null;
    }
}
function infinite_loop_guard() {
    try {
        _errors_js__WEBPACK_IMPORTED_MODULE_14__/* .effect_update_depth_exceeded */.Cl();
    } catch (error) {
        if (esm_env__WEBPACK_IMPORTED_MODULE_15__/* ["default"] */.A) {
            // stack contains no useful information, replace it
            (0,_shared_utils_js__WEBPACK_IMPORTED_MODULE_1__/* .define_property */.Qu)(error, 'stack', {
                value: ''
            });
        }
        // Best effort: invoke the boundary nearest the most recent
        // effect and hope that it's relevant to the infinite loop
        (0,_error_handling_js__WEBPACK_IMPORTED_MODULE_4__/* .invoke_error_boundary */.n)(error, last_scheduled_effect);
    }
}
/** @type {Effect[] | null} */ let eager_block_effects = null;
/**
 * @param {Array<Effect>} effects
 * @returns {void}
 */ function flush_queued_effects(effects) {
    var length = effects.length;
    if (length === 0) return;
    var i = 0;
    while(i < length){
        var effect = effects[i++];
        if ((effect.f & (_client_constants__WEBPACK_IMPORTED_MODULE_0__/* .DESTROYED */.o5 | _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .INERT */.$q)) === 0 && (0,_runtime_js__WEBPACK_IMPORTED_MODULE_2__/* .is_dirty */.Kj)(effect)) {
            eager_block_effects = [];
            (0,_runtime_js__WEBPACK_IMPORTED_MODULE_2__/* .update_effect */.gJ)(effect);
            // Effects with no dependencies or teardown do not get added to the effect tree.
            // Deferred effects (e.g. `$effect(...)`) _are_ added to the tree because we
            // don't know if we need to keep them until they are executed. Doing the check
            // here (rather than in `update_effect`) allows us to skip the work for
            // immediate effects.
            if (effect.deps === null && effect.first === null && effect.nodes_start === null) {
                // if there's no teardown or abort controller we completely unlink
                // the effect from the graph
                if (effect.teardown === null && effect.ac === null) {
                    // remove this effect from the graph
                    (0,_effects_js__WEBPACK_IMPORTED_MODULE_6__/* .unlink_effect */.qX)(effect);
                } else {
                    // keep the effect in the graph, but free up some memory
                    effect.fn = null;
                }
            }
            // If update_effect() has a flushSync() in it, we may have flushed another flush_queued_effects(),
            // which already handled this logic and did set eager_block_effects to null.
            if ((eager_block_effects === null || eager_block_effects === void 0 ? void 0 : eager_block_effects.length) > 0) {
                // TODO this feels incorrect! it gets the tests passing
                _sources_js__WEBPACK_IMPORTED_MODULE_5__/* .old_values.clear */.bJ.clear();
                for (const e of eager_block_effects){
                    (0,_runtime_js__WEBPACK_IMPORTED_MODULE_2__/* .update_effect */.gJ)(e);
                }
                eager_block_effects = [];
            }
        }
    }
    eager_block_effects = null;
}
/**
 * @param {Effect} signal
 * @returns {void}
 */ function schedule_effect(signal) {
    var effect = last_scheduled_effect = signal;
    while(effect.parent !== null){
        effect = effect.parent;
        var flags = effect.f;
        // if the effect is being scheduled because a parent (each/await/etc) block
        // updated an internal source, bail out or we'll cause a second flush
        if (is_flushing && effect === _runtime_js__WEBPACK_IMPORTED_MODULE_2__/* .active_effect */.Fg && (flags & _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .BLOCK_EFFECT */.kc) !== 0) {
            return;
        }
        if ((flags & (_client_constants__WEBPACK_IMPORTED_MODULE_0__/* .ROOT_EFFECT */.FV | _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .BRANCH_EFFECT */.Zr)) !== 0) {
            if ((flags & _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .CLEAN */.w_) === 0) return;
            effect.f ^= _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .CLEAN */.w_;
        }
    }
    queued_root_effects.push(effect);
}
/**
 * Forcibly remove all current batches, to prevent cross-talk between tests
 */ function clear() {
    batches.clear();
}


}),
462: (function (__unused_webpack_module, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  V: () => (set_from_async_derived),
  Xd: () => (derived_safe_equal),
  c2: () => (update_derived),
  eO: () => (user_derived),
  ge: () => (destroy_derived_effects),
  kX: () => (recent_async_deriveds),
  un: () => (derived),
  vO: () => (current_async_effect),
  w6: () => (execute_derived),
  zx: () => (async_derived)
});
/* ESM import */var esm_env__WEBPACK_IMPORTED_MODULE_10__ = __webpack_require__(832);
/* ESM import */var _client_constants__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(924);
/* ESM import */var _runtime_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(513);
/* ESM import */var _equality_js__WEBPACK_IMPORTED_MODULE_9__ = __webpack_require__(576);
/* ESM import */var _errors_js__WEBPACK_IMPORTED_MODULE_12__ = __webpack_require__(626);
/* ESM import */var _warnings_js__WEBPACK_IMPORTED_MODULE_13__ = __webpack_require__(32);
/* ESM import */var _effects_js__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(480);
/* ESM import */var _sources_js__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(264);
/* ESM import */var _dev_tracing_js__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(339);
/* ESM import */var _flags_index_js__WEBPACK_IMPORTED_MODULE_11__ = __webpack_require__(817);
/* ESM import */var _context_js__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(754);
/* ESM import */var _constants_js__WEBPACK_IMPORTED_MODULE_6__ = __webpack_require__(178);
/* ESM import */var _batch_js__WEBPACK_IMPORTED_MODULE_7__ = __webpack_require__(410);
/* ESM import */var _async_js__WEBPACK_IMPORTED_MODULE_8__ = __webpack_require__(850);
/** @import { Derived, Effect, Source } from '#client' */ /** @import { Batch } from './batch.js'; */ 













/** @type {Effect | null} */ let current_async_effect = null;
/** @param {Effect | null} v */ function set_from_async_derived(v) {
    current_async_effect = v;
}
const recent_async_deriveds = new Set();
/**
 * @template V
 * @param {() => V} fn
 * @returns {Derived<V>}
 */ /*#__NO_SIDE_EFFECTS__*/ function derived(fn) {
    var flags = _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .DERIVED */.mj | _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .DIRTY */.jm;
    var parent_derived = _runtime_js__WEBPACK_IMPORTED_MODULE_1__/* .active_reaction */.hp !== null && (_runtime_js__WEBPACK_IMPORTED_MODULE_1__/* .active_reaction.f */.hp.f & _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .DERIVED */.mj) !== 0 ? /** @type {Derived} */ _runtime_js__WEBPACK_IMPORTED_MODULE_1__/* .active_reaction */.hp : null;
    if (_runtime_js__WEBPACK_IMPORTED_MODULE_1__/* .active_effect */.Fg === null || parent_derived !== null && (parent_derived.f & _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .UNOWNED */.L2) !== 0) {
        flags |= _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .UNOWNED */.L2;
    } else {
        // Since deriveds are evaluated lazily, any effects created inside them are
        // created too late to ensure that the parent effect is added to the tree
        _runtime_js__WEBPACK_IMPORTED_MODULE_1__/* .active_effect.f */.Fg.f |= _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .EFFECT_PRESERVED */.V$;
    }
    /** @type {Derived<V>} */ const signal = {
        ctx: _context_js__WEBPACK_IMPORTED_MODULE_5__/* .component_context */.UL,
        deps: null,
        effects: null,
        equals: _equality_js__WEBPACK_IMPORTED_MODULE_9__/* .equals */.aI,
        f: flags,
        fn,
        reactions: null,
        rv: 0,
        v: /** @type {V} */ _constants_js__WEBPACK_IMPORTED_MODULE_6__/* .UNINITIALIZED */.UP,
        wv: 0,
        parent: parent_derived ?? _runtime_js__WEBPACK_IMPORTED_MODULE_1__/* .active_effect */.Fg,
        ac: null
    };
    if (esm_env__WEBPACK_IMPORTED_MODULE_10__/* ["default"] */.A && _flags_index_js__WEBPACK_IMPORTED_MODULE_11__/* .tracing_mode_flag */._G) {
        signal.created = (0,_dev_tracing_js__WEBPACK_IMPORTED_MODULE_4__/* .get_stack */.sv)('CreatedAt');
    }
    return signal;
}
/**
 * @template V
 * @param {() => V | Promise<V>} fn
 * @param {string} [location] If provided, print a warning if the value is not read immediately after update
 * @returns {Promise<Source<V>>}
 */ /*#__NO_SIDE_EFFECTS__*/ function async_derived(fn, location) {
    let parent = /** @type {Effect | null} */ _runtime_js__WEBPACK_IMPORTED_MODULE_1__/* .active_effect */.Fg;
    if (parent === null) {
        _errors_js__WEBPACK_IMPORTED_MODULE_12__/* .async_derived_orphan */.aQ();
    }
    var boundary = /** @type {Boundary} */ parent.b;
    var promise = /** @type {unknown} */ undefined;
    var signal = (0,_sources_js__WEBPACK_IMPORTED_MODULE_3__/* .source */.sP)(/** @type {V} */ _constants_js__WEBPACK_IMPORTED_MODULE_6__/* .UNINITIALIZED */.UP);
    /** @type {Promise<V> | null} */ var prev = null;
    // only suspend in async deriveds created on initialisation
    var should_suspend = !_runtime_js__WEBPACK_IMPORTED_MODULE_1__/* .active_reaction */.hp;
    (0,_effects_js__WEBPACK_IMPORTED_MODULE_2__/* .async_effect */.NQ)(()=>{
        if (esm_env__WEBPACK_IMPORTED_MODULE_10__/* ["default"] */.A) current_async_effect = _runtime_js__WEBPACK_IMPORTED_MODULE_1__/* .active_effect */.Fg;
        try {
            var p = fn();
            // Make sure to always access the then property to read any signals
            // it might access, so that we track them as dependencies.
            if (prev) Promise.resolve(p).catch(()=>{}); // avoid unhandled rejection
        } catch (error) {
            p = Promise.reject(error);
        }
        if (esm_env__WEBPACK_IMPORTED_MODULE_10__/* ["default"] */.A) current_async_effect = null;
        var r = ()=>p;
        promise = (prev === null || prev === void 0 ? void 0 : prev.then(r, r)) ?? Promise.resolve(p);
        prev = promise;
        var batch = /** @type {Batch} */ _batch_js__WEBPACK_IMPORTED_MODULE_7__/* .current_batch */.Dr;
        var pending = boundary.is_pending();
        if (should_suspend) {
            boundary.update_pending_count(1);
            if (!pending) batch.increment();
        }
        /**
		 * @param {any} value
		 * @param {unknown} error
		 */ const handler = function(value) {
            let error = arguments.length > 1 && arguments[1] !== void 0 ? arguments[1] : undefined;
            prev = null;
            current_async_effect = null;
            if (!pending) batch.activate();
            if (error) {
                if (error !== _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .STALE_REACTION */.In) {
                    signal.f |= _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .ERROR_VALUE */.dH;
                    // @ts-expect-error the error is the wrong type, but we don't care
                    (0,_sources_js__WEBPACK_IMPORTED_MODULE_3__/* .internal_set */.LY)(signal, error);
                }
            } else {
                if ((signal.f & _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .ERROR_VALUE */.dH) !== 0) {
                    signal.f ^= _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .ERROR_VALUE */.dH;
                }
                (0,_sources_js__WEBPACK_IMPORTED_MODULE_3__/* .internal_set */.LY)(signal, value);
                if (esm_env__WEBPACK_IMPORTED_MODULE_10__/* ["default"] */.A && location !== undefined) {
                    recent_async_deriveds.add(signal);
                    setTimeout(()=>{
                        if (recent_async_deriveds.has(signal)) {
                            _warnings_js__WEBPACK_IMPORTED_MODULE_13__/* .await_waterfall */.Cy(/** @type {string} */ signal.label, location);
                            recent_async_deriveds.delete(signal);
                        }
                    });
                }
            }
            if (should_suspend) {
                boundary.update_pending_count(-1);
                if (!pending) batch.decrement();
            }
            (0,_async_js__WEBPACK_IMPORTED_MODULE_8__/* .unset_context */.sO)();
        };
        promise.then(handler, (e)=>handler(null, e || 'unknown'));
        if (batch) {
            return ()=>{
                queueMicrotask(()=>batch.neuter());
            };
        }
    });
    if (esm_env__WEBPACK_IMPORTED_MODULE_10__/* ["default"] */.A) {
        // add a flag that lets this be printed as a derived
        // when using `$inspect.trace()`
        signal.f |= _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .ASYNC */.VD;
    }
    return new Promise((fulfil)=>{
        /** @param {Promise<V>} p */ function next(p) {
            function go() {
                if (p === promise) {
                    fulfil(signal);
                } else {
                    // if the effect re-runs before the initial promise
                    // resolves, delay resolution until we have a value
                    next(promise);
                }
            }
            p.then(go, go);
        }
        next(promise);
    });
}
/**
 * @template V
 * @param {() => V} fn
 * @returns {Derived<V>}
 */ /*#__NO_SIDE_EFFECTS__*/ function user_derived(fn) {
    const d = derived(fn);
    (0,_runtime_js__WEBPACK_IMPORTED_MODULE_1__/* .push_reaction_value */.tT)(d);
    return d;
}
/**
 * @template V
 * @param {() => V} fn
 * @returns {Derived<V>}
 */ /*#__NO_SIDE_EFFECTS__*/ function derived_safe_equal(fn) {
    const signal = derived(fn);
    signal.equals = _equality_js__WEBPACK_IMPORTED_MODULE_9__/* .safe_equals */.Og;
    return signal;
}
/**
 * @param {Derived} derived
 * @returns {void}
 */ function destroy_derived_effects(derived) {
    var effects = derived.effects;
    if (effects !== null) {
        derived.effects = null;
        for(var i = 0; i < effects.length; i += 1){
            (0,_effects_js__WEBPACK_IMPORTED_MODULE_2__/* .destroy_effect */.DI)(/** @type {Effect} */ effects[i]);
        }
    }
}
/**
 * The currently updating deriveds, used to detect infinite recursion
 * in dev mode and provide a nicer error than 'too much recursion'
 * @type {Derived[]}
 */ let stack = [];
/**
 * @param {Derived} derived
 * @returns {Effect | null}
 */ function get_derived_parent_effect(derived) {
    var parent = derived.parent;
    while(parent !== null){
        if ((parent.f & _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .DERIVED */.mj) === 0) {
            return /** @type {Effect} */ parent;
        }
        parent = parent.parent;
    }
    return null;
}
/**
 * @template T
 * @param {Derived} derived
 * @returns {T}
 */ function execute_derived(derived) {
    var value;
    var prev_active_effect = _runtime_js__WEBPACK_IMPORTED_MODULE_1__/* .active_effect */.Fg;
    (0,_runtime_js__WEBPACK_IMPORTED_MODULE_1__/* .set_active_effect */.gU)(get_derived_parent_effect(derived));
    if (esm_env__WEBPACK_IMPORTED_MODULE_10__/* ["default"] */.A) {
        let prev_inspect_effects = _sources_js__WEBPACK_IMPORTED_MODULE_3__/* .inspect_effects */.MU;
        (0,_sources_js__WEBPACK_IMPORTED_MODULE_3__/* .set_inspect_effects */.JY)(new Set());
        try {
            if (stack.includes(derived)) {
                _errors_js__WEBPACK_IMPORTED_MODULE_12__/* .derived_references_self */.cN();
            }
            stack.push(derived);
            destroy_derived_effects(derived);
            value = (0,_runtime_js__WEBPACK_IMPORTED_MODULE_1__/* .update_reaction */.mj)(derived);
        } finally{
            (0,_runtime_js__WEBPACK_IMPORTED_MODULE_1__/* .set_active_effect */.gU)(prev_active_effect);
            (0,_sources_js__WEBPACK_IMPORTED_MODULE_3__/* .set_inspect_effects */.JY)(prev_inspect_effects);
            stack.pop();
        }
    } else {
        try {
            destroy_derived_effects(derived);
            value = (0,_runtime_js__WEBPACK_IMPORTED_MODULE_1__/* .update_reaction */.mj)(derived);
        } finally{
            (0,_runtime_js__WEBPACK_IMPORTED_MODULE_1__/* .set_active_effect */.gU)(prev_active_effect);
        }
    }
    return value;
}
/**
 * @param {Derived} derived
 * @returns {void}
 */ function update_derived(derived) {
    var value = execute_derived(derived);
    if (!derived.equals(value)) {
        derived.v = value;
        derived.wv = (0,_runtime_js__WEBPACK_IMPORTED_MODULE_1__/* .increment_write_version */.Fq)();
    }
    // don't mark derived clean if we're reading it inside a
    // cleanup function, or it will cache a stale value
    if (_runtime_js__WEBPACK_IMPORTED_MODULE_1__/* .is_destroying_effect */.WI) {
        return;
    }
    if (_batch_js__WEBPACK_IMPORTED_MODULE_7__/* .batch_deriveds */.G1 !== null) {
        _batch_js__WEBPACK_IMPORTED_MODULE_7__/* .batch_deriveds.set */.G1.set(derived, derived.v);
    } else {
        var status = (_runtime_js__WEBPACK_IMPORTED_MODULE_1__/* .skip_reaction */.U9 || (derived.f & _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .UNOWNED */.L2) !== 0) && derived.deps !== null ? _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .MAYBE_DIRTY */.ig : _client_constants__WEBPACK_IMPORTED_MODULE_0__/* .CLEAN */.w_;
        (0,_runtime_js__WEBPACK_IMPORTED_MODULE_1__/* .set_signal_status */.TC)(derived, status);
    }
}


}),
480: (function (__unused_webpack_module, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  DI: () => (destroy_effect),
  F3: () => (destroy_effect_children),
  Fc: () => (effect_root),
  MW: () => (user_effect),
  NQ: () => (async_effect),
  Nq: () => (execute_effect_teardown),
  QZ: () => (effect),
  V1: () => (create_user_effect),
  VB: () => (render_effect),
  cc: () => (resume_effect),
  mk: () => (remove_effect_dom),
  oJ: () => (effect_tracking),
  om: () => (block),
  pk: () => (destroy_block_effect_children),
  qX: () => (unlink_effect),
  r4: () => (pause_effect),
  tk: () => (branch),
  vN: () => (template_effect),
  x4: () => (component_root)
});
/* ESM import */var _runtime_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(513);
/* ESM import */var _client_constants__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(924);
/* ESM import */var _errors_js__WEBPACK_IMPORTED_MODULE_8__ = __webpack_require__(626);
/* ESM import */var esm_env__WEBPACK_IMPORTED_MODULE_9__ = __webpack_require__(832);
/* ESM import */var _shared_utils_js__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(986);
/* ESM import */var _dom_operations_js__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(518);
/* ESM import */var _context_js__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(754);
/* ESM import */var _batch_js__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(410);
/* ESM import */var _async_js__WEBPACK_IMPORTED_MODULE_6__ = __webpack_require__(850);
/* ESM import */var _dom_elements_bindings_shared_js__WEBPACK_IMPORTED_MODULE_7__ = __webpack_require__(408);
/** @import { ComponentContext, ComponentContextLegacy, Derived, Effect, TemplateNode, TransitionManager } from '#client' */ 









/**
 * @param {'$effect' | '$effect.pre' | '$inspect'} rune
 */ function validate_effect(rune) {
    if (_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .active_effect */.Fg === null && _runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .active_reaction */.hp === null) {
        _errors_js__WEBPACK_IMPORTED_MODULE_8__/* .effect_orphan */.tB(rune);
    }
    if (_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .active_reaction */.hp !== null && (_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .active_reaction.f */.hp.f & _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .UNOWNED */.L2) !== 0 && _runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .active_effect */.Fg === null) {
        _errors_js__WEBPACK_IMPORTED_MODULE_8__/* .effect_in_unowned_derived */.fi();
    }
    if (_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .is_destroying_effect */.WI) {
        _errors_js__WEBPACK_IMPORTED_MODULE_8__/* .effect_in_teardown */.BT(rune);
    }
}
/**
 * @param {Effect} effect
 * @param {Effect} parent_effect
 */ function push_effect(effect, parent_effect) {
    var parent_last = parent_effect.last;
    if (parent_last === null) {
        parent_effect.last = parent_effect.first = effect;
    } else {
        parent_last.next = effect;
        effect.prev = parent_last;
        parent_effect.last = effect;
    }
}
/**
 * @param {number} type
 * @param {null | (() => void | (() => void))} fn
 * @param {boolean} sync
 * @param {boolean} push
 * @returns {Effect}
 */ function create_effect(type, fn, sync) {
    let push = arguments.length > 3 && arguments[3] !== void 0 ? arguments[3] : true;
    var parent = _runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .active_effect */.Fg;
    if (esm_env__WEBPACK_IMPORTED_MODULE_9__/* ["default"] */.A) {
        // Ensure the parent is never an inspect effect
        while(parent !== null && (parent.f & _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .INSPECT_EFFECT */.T1) !== 0){
            parent = parent.parent;
        }
    }
    if (parent !== null && (parent.f & _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .INERT */.$q) !== 0) {
        type |= _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .INERT */.$q;
    }
    /** @type {Effect} */ var effect = {
        ctx: _context_js__WEBPACK_IMPORTED_MODULE_4__/* .component_context */.UL,
        deps: null,
        nodes_start: null,
        nodes_end: null,
        f: type | _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .DIRTY */.jm,
        first: null,
        fn,
        last: null,
        next: null,
        parent,
        b: parent && parent.b,
        prev: null,
        teardown: null,
        transitions: null,
        wv: 0,
        ac: null
    };
    if (esm_env__WEBPACK_IMPORTED_MODULE_9__/* ["default"] */.A) {
        effect.component_function = _context_js__WEBPACK_IMPORTED_MODULE_4__/* .dev_current_component_function */.DE;
    }
    if (sync) {
        try {
            (0,_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .update_effect */.gJ)(effect);
            effect.f |= _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .EFFECT_RAN */.wi;
        } catch (e) {
            destroy_effect(effect);
            throw e;
        }
    } else if (fn !== null) {
        (0,_batch_js__WEBPACK_IMPORTED_MODULE_5__/* .schedule_effect */.ec)(effect);
    }
    if (push) {
        /** @type {Effect | null} */ var e = effect;
        // if an effect has already ran and doesn't need to be kept in the tree
        // (because it won't re-run, has no DOM, and has no teardown etc)
        // then we skip it and go to its child (if any)
        if (sync && e.deps === null && e.teardown === null && e.nodes_start === null && e.first === e.last && // either `null`, or a singular child
        (e.f & _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .EFFECT_PRESERVED */.V$) === 0) {
            e = e.first;
        }
        if (e !== null) {
            e.parent = parent;
            if (parent !== null) {
                push_effect(e, parent);
            }
            // if we're in a derived, add the effect there too
            if (_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .active_reaction */.hp !== null && (_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .active_reaction.f */.hp.f & _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .DERIVED */.mj) !== 0 && (type & _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .ROOT_EFFECT */.FV) === 0) {
                var _derived;
                var derived = /** @type {Derived} */ _runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .active_reaction */.hp;
                ((_derived = derived).effects ?? (_derived.effects = [])).push(e);
            }
        }
    }
    return effect;
}
/**
 * Internal representation of `$effect.tracking()`
 * @returns {boolean}
 */ function effect_tracking() {
    return _runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .active_reaction */.hp !== null && !_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .untracking */.LW;
}
/**
 * @param {() => void} fn
 */ function teardown(fn) {
    const effect = create_effect(RENDER_EFFECT, null, false);
    set_signal_status(effect, CLEAN);
    effect.teardown = fn;
    return effect;
}
/**
 * Internal representation of `$effect(...)`
 * @param {() => void | (() => void)} fn
 */ function user_effect(fn) {
    validate_effect('$effect');
    if (esm_env__WEBPACK_IMPORTED_MODULE_9__/* ["default"] */.A) {
        (0,_shared_utils_js__WEBPACK_IMPORTED_MODULE_2__/* .define_property */.Qu)(fn, 'name', {
            value: '$effect'
        });
    }
    // Non-nested `$effect(...)` in a component should be deferred
    // until the component is mounted
    var flags = /** @type {Effect} */ _runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .active_effect.f */.Fg.f;
    var defer = !_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .active_reaction */.hp && (flags & _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .BRANCH_EFFECT */.Zr) !== 0 && (flags & _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .EFFECT_RAN */.wi) === 0;
    if (defer) {
        var _context;
        // Top-level `$effect(...)` in an unmounted component — defer until mount
        var context = /** @type {ComponentContext} */ _context_js__WEBPACK_IMPORTED_MODULE_4__/* .component_context */.UL;
        ((_context = context).e ?? (_context.e = [])).push(fn);
    } else {
        // Everything else — create immediately
        return create_user_effect(fn);
    }
}
/**
 * @param {() => void | (() => void)} fn
 */ function create_user_effect(fn) {
    return create_effect(_client_constants__WEBPACK_IMPORTED_MODULE_1__/* .EFFECT */.ac | _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .USER_EFFECT */.Wr, fn, false);
}
/**
 * Internal representation of `$effect.pre(...)`
 * @param {() => void | (() => void)} fn
 * @returns {Effect}
 */ function user_pre_effect(fn) {
    validate_effect('$effect.pre');
    if (DEV) {
        define_property(fn, 'name', {
            value: '$effect.pre'
        });
    }
    return create_effect(RENDER_EFFECT | USER_EFFECT, fn, true);
}
/** @param {() => void | (() => void)} fn */ function inspect_effect(fn) {
    return create_effect(INSPECT_EFFECT, fn, true);
}
/**
 * Internal representation of `$effect.root(...)`
 * @param {() => void | (() => void)} fn
 * @returns {() => void}
 */ function effect_root(fn) {
    _batch_js__WEBPACK_IMPORTED_MODULE_5__/* .Batch.ensure */.lP.ensure();
    const effect = create_effect(_client_constants__WEBPACK_IMPORTED_MODULE_1__/* .ROOT_EFFECT */.FV | _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .EFFECT_PRESERVED */.V$, fn, true);
    return ()=>{
        destroy_effect(effect);
    };
}
/**
 * An effect root whose children can transition out
 * @param {() => void} fn
 * @returns {(options?: { outro?: boolean }) => Promise<void>}
 */ function component_root(fn) {
    _batch_js__WEBPACK_IMPORTED_MODULE_5__/* .Batch.ensure */.lP.ensure();
    const effect = create_effect(_client_constants__WEBPACK_IMPORTED_MODULE_1__/* .ROOT_EFFECT */.FV | _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .EFFECT_PRESERVED */.V$, fn, true);
    return function() {
        let options = arguments.length > 0 && arguments[0] !== void 0 ? arguments[0] : {};
        return new Promise((fulfil)=>{
            if (options.outro) {
                pause_effect(effect, ()=>{
                    destroy_effect(effect);
                    fulfil(undefined);
                });
            } else {
                destroy_effect(effect);
                fulfil(undefined);
            }
        });
    };
}
/**
 * @param {() => void | (() => void)} fn
 * @returns {Effect}
 */ function effect(fn) {
    return create_effect(_client_constants__WEBPACK_IMPORTED_MODULE_1__/* .EFFECT */.ac, fn, false);
}
/**
 * Internal representation of `$: ..`
 * @param {() => any} deps
 * @param {() => void | (() => void)} fn
 */ function legacy_pre_effect(deps, fn) {
    var context = /** @type {ComponentContextLegacy} */ component_context;
    /** @type {{ effect: null | Effect, ran: boolean, deps: () => any }} */ var token = {
        effect: null,
        ran: false,
        deps
    };
    context.l.$.push(token);
    token.effect = render_effect(()=>{
        deps();
        // If this legacy pre effect has already run before the end of the reset, then
        // bail out to emulate the same behavior.
        if (token.ran) return;
        token.ran = true;
        untrack(fn);
    });
}
function legacy_pre_effect_reset() {
    var context = /** @type {ComponentContextLegacy} */ component_context;
    render_effect(()=>{
        // Run dirty `$:` statements
        for (var token of context.l.$){
            token.deps();
            var effect = token.effect;
            // If the effect is CLEAN, then make it MAYBE_DIRTY. This ensures we traverse through
            // the effects dependencies and correctly ensure each dependency is up-to-date.
            if ((effect.f & CLEAN) !== 0) {
                set_signal_status(effect, MAYBE_DIRTY);
            }
            if (is_dirty(effect)) {
                update_effect(effect);
            }
            token.ran = false;
        }
    });
}
/**
 * @param {() => void | (() => void)} fn
 * @returns {Effect}
 */ function async_effect(fn) {
    return create_effect(_client_constants__WEBPACK_IMPORTED_MODULE_1__/* .ASYNC */.VD | _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .EFFECT_PRESERVED */.V$, fn, true);
}
/**
 * @param {() => void | (() => void)} fn
 * @returns {Effect}
 */ function render_effect(fn) {
    let flags = arguments.length > 1 && arguments[1] !== void 0 ? arguments[1] : 0;
    return create_effect(_client_constants__WEBPACK_IMPORTED_MODULE_1__/* .RENDER_EFFECT */.Zv | flags, fn, true);
}
/**
 * @param {(...expressions: any) => void | (() => void)} fn
 * @param {Array<() => any>} sync
 * @param {Array<() => Promise<any>>} async
 */ function template_effect(fn) {
    let sync = arguments.length > 1 && arguments[1] !== void 0 ? arguments[1] : [], async = arguments.length > 2 && arguments[2] !== void 0 ? arguments[2] : [];
    (0,_async_js__WEBPACK_IMPORTED_MODULE_6__/* .flatten */.Bq)(sync, async, (values)=>{
        create_effect(_client_constants__WEBPACK_IMPORTED_MODULE_1__/* .RENDER_EFFECT */.Zv, ()=>fn(...values.map(_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .get */.Jt)), true);
    });
}
/**
 * @param {(() => void)} fn
 * @param {number} flags
 */ function block(fn) {
    let flags = arguments.length > 1 && arguments[1] !== void 0 ? arguments[1] : 0;
    var effect = create_effect(_client_constants__WEBPACK_IMPORTED_MODULE_1__/* .BLOCK_EFFECT */.kc | flags, fn, true);
    if (esm_env__WEBPACK_IMPORTED_MODULE_9__/* ["default"] */.A) {
        effect.dev_stack = _context_js__WEBPACK_IMPORTED_MODULE_4__/* .dev_stack */.lv;
    }
    return effect;
}
/**
 * @param {(() => void)} fn
 * @param {boolean} [push]
 */ function branch(fn) {
    let push = arguments.length > 1 && arguments[1] !== void 0 ? arguments[1] : true;
    return create_effect(_client_constants__WEBPACK_IMPORTED_MODULE_1__/* .BRANCH_EFFECT */.Zr | _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .EFFECT_PRESERVED */.V$, fn, true, push);
}
/**
 * @param {Effect} effect
 */ function execute_effect_teardown(effect) {
    var teardown = effect.teardown;
    if (teardown !== null) {
        const previously_destroying_effect = _runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .is_destroying_effect */.WI;
        const previous_reaction = _runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .active_reaction */.hp;
        (0,_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .set_is_destroying_effect */.fT)(true);
        (0,_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .set_active_reaction */.G0)(null);
        try {
            teardown.call(null);
        } finally{
            (0,_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .set_is_destroying_effect */.fT)(previously_destroying_effect);
            (0,_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .set_active_reaction */.G0)(previous_reaction);
        }
    }
}
/**
 * @param {Effect} signal
 * @param {boolean} remove_dom
 * @returns {void}
 */ function destroy_effect_children(signal) {
    let remove_dom = arguments.length > 1 && arguments[1] !== void 0 ? arguments[1] : false;
    var effect = signal.first;
    signal.first = signal.last = null;
    while(effect !== null){
        const controller = effect.ac;
        if (controller !== null) {
            (0,_dom_elements_bindings_shared_js__WEBPACK_IMPORTED_MODULE_7__/* .without_reactive_context */.$w)(()=>{
                controller.abort(_client_constants__WEBPACK_IMPORTED_MODULE_1__/* .STALE_REACTION */.In);
            });
        }
        var next = effect.next;
        if ((effect.f & _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .ROOT_EFFECT */.FV) !== 0) {
            // this is now an independent root
            effect.parent = null;
        } else {
            destroy_effect(effect, remove_dom);
        }
        effect = next;
    }
}
/**
 * @param {Effect} signal
 * @returns {void}
 */ function destroy_block_effect_children(signal) {
    var effect = signal.first;
    while(effect !== null){
        var next = effect.next;
        if ((effect.f & _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .BRANCH_EFFECT */.Zr) === 0) {
            destroy_effect(effect);
        }
        effect = next;
    }
}
/**
 * @param {Effect} effect
 * @param {boolean} [remove_dom]
 * @returns {void}
 */ function destroy_effect(effect) {
    let remove_dom = arguments.length > 1 && arguments[1] !== void 0 ? arguments[1] : true;
    var removed = false;
    if ((remove_dom || (effect.f & _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .HEAD_EFFECT */.PL) !== 0) && effect.nodes_start !== null && effect.nodes_end !== null) {
        remove_effect_dom(effect.nodes_start, /** @type {TemplateNode} */ effect.nodes_end);
        removed = true;
    }
    destroy_effect_children(effect, remove_dom && !removed);
    (0,_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .remove_reactions */.yR)(effect, 0);
    (0,_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .set_signal_status */.TC)(effect, _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .DESTROYED */.o5);
    var transitions = effect.transitions;
    if (transitions !== null) {
        for (const transition of transitions){
            transition.stop();
        }
    }
    execute_effect_teardown(effect);
    var parent = effect.parent;
    // If the parent doesn't have any children, then skip this work altogether
    if (parent !== null && parent.first !== null) {
        unlink_effect(effect);
    }
    if (esm_env__WEBPACK_IMPORTED_MODULE_9__/* ["default"] */.A) {
        effect.component_function = null;
    }
    // `first` and `child` are nulled out in destroy_effect_children
    // we don't null out `parent` so that error propagation can work correctly
    effect.next = effect.prev = effect.teardown = effect.ctx = effect.deps = effect.fn = effect.nodes_start = effect.nodes_end = effect.ac = null;
}
/**
 *
 * @param {TemplateNode | null} node
 * @param {TemplateNode} end
 */ function remove_effect_dom(node, end) {
    while(node !== null){
        /** @type {TemplateNode | null} */ var next = node === end ? null : /** @type {TemplateNode} */ (0,_dom_operations_js__WEBPACK_IMPORTED_MODULE_3__/* .get_next_sibling */.M$)(node);
        node.remove();
        node = next;
    }
}
/**
 * Detach an effect from the effect tree, freeing up memory and
 * reducing the amount of work that happens on subsequent traversals
 * @param {Effect} effect
 */ function unlink_effect(effect) {
    var parent = effect.parent;
    var prev = effect.prev;
    var next = effect.next;
    if (prev !== null) prev.next = next;
    if (next !== null) next.prev = prev;
    if (parent !== null) {
        if (parent.first === effect) parent.first = next;
        if (parent.last === effect) parent.last = prev;
    }
}
/**
 * When a block effect is removed, we don't immediately destroy it or yank it
 * out of the DOM, because it might have transitions. Instead, we 'pause' it.
 * It stays around (in memory, and in the DOM) until outro transitions have
 * completed, and if the state change is reversed then we _resume_ it.
 * A paused effect does not update, and the DOM subtree becomes inert.
 * @param {Effect} effect
 * @param {() => void} [callback]
 */ function pause_effect(effect, callback) {
    /** @type {TransitionManager[]} */ var transitions = [];
    pause_children(effect, transitions, true);
    run_out_transitions(transitions, ()=>{
        destroy_effect(effect);
        if (callback) callback();
    });
}
/**
 * @param {TransitionManager[]} transitions
 * @param {() => void} fn
 */ function run_out_transitions(transitions, fn) {
    var remaining = transitions.length;
    if (remaining > 0) {
        var check = ()=>--remaining || fn();
        for (var transition of transitions){
            transition.out(check);
        }
    } else {
        fn();
    }
}
/**
 * @param {Effect} effect
 * @param {TransitionManager[]} transitions
 * @param {boolean} local
 */ function pause_children(effect, transitions, local) {
    if ((effect.f & _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .INERT */.$q) !== 0) return;
    effect.f ^= _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .INERT */.$q;
    if (effect.transitions !== null) {
        for (const transition of effect.transitions){
            if (transition.is_global || local) {
                transitions.push(transition);
            }
        }
    }
    var child = effect.first;
    while(child !== null){
        var sibling = child.next;
        var transparent = (child.f & _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .EFFECT_TRANSPARENT */.lQ) !== 0 || (child.f & _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .BRANCH_EFFECT */.Zr) !== 0;
        // TODO we don't need to call pause_children recursively with a linked list in place
        // it's slightly more involved though as we have to account for `transparent` changing
        // through the tree.
        pause_children(child, transitions, transparent ? local : false);
        child = sibling;
    }
}
/**
 * The opposite of `pause_effect`. We call this if (for example)
 * `x` becomes falsy then truthy: `{#if x}...{/if}`
 * @param {Effect} effect
 */ function resume_effect(effect) {
    resume_children(effect, true);
}
/**
 * @param {Effect} effect
 * @param {boolean} local
 */ function resume_children(effect, local) {
    if ((effect.f & _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .INERT */.$q) === 0) return;
    effect.f ^= _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .INERT */.$q;
    // If a dependency of this effect changed while it was paused,
    // schedule the effect to update. we don't use `is_dirty`
    // here because we don't want to eagerly recompute a derived like
    // `{#if foo}{foo.bar()}{/if}` if `foo` is now `undefined
    if ((effect.f & _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .CLEAN */.w_) === 0) {
        (0,_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .set_signal_status */.TC)(effect, _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .DIRTY */.jm);
        (0,_batch_js__WEBPACK_IMPORTED_MODULE_5__/* .schedule_effect */.ec)(effect);
    }
    var child = effect.first;
    while(child !== null){
        var sibling = child.next;
        var transparent = (child.f & _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .EFFECT_TRANSPARENT */.lQ) !== 0 || (child.f & _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .BRANCH_EFFECT */.Zr) !== 0;
        // TODO we don't need to call resume_children recursively with a linked list in place
        // it's slightly more involved though as we have to account for `transparent` changing
        // through the tree.
        resume_children(child, transparent ? local : false);
        child = sibling;
    }
    if (effect.transitions !== null) {
        for (const transition of effect.transitions){
            if (transition.is_global || local) {
                transition.in();
            }
        }
    }
}
function aborted() {
    let effect = arguments.length > 0 && arguments[0] !== void 0 ? arguments[0] : /** @type {Effect} */ active_effect;
    return (effect.f & DESTROYED) !== 0;
}


}),
576: (function (__unused_webpack_module, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  Og: () => (safe_equals),
  aI: () => (equals),
  jX: () => (safe_not_equal)
});
/** @import { Equals } from '#client' */ /** @type {Equals} */ function equals(value) {
    return value === this.v;
}
/**
 * @param {unknown} a
 * @param {unknown} b
 * @returns {boolean}
 */ function safe_not_equal(a, b) {
    return a != a ? b == b : a !== b || a !== null && typeof a === 'object' || typeof a === 'function';
}
/**
 * @param {unknown} a
 * @param {unknown} b
 * @returns {boolean}
 */ function not_equal(a, b) {
    return a !== b;
}
/** @type {Equals} */ function safe_equals(value) {
    return !safe_not_equal(value, this.v);
}


}),
264: (function (__unused_webpack_module, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  GV: () => (increment),
  HX: () => (set_inspect_effects_deferred),
  JY: () => (set_inspect_effects),
  LY: () => (internal_set),
  MU: () => (inspect_effects),
  Xy: () => (flush_inspect_effects),
  bJ: () => (old_values),
  hZ: () => (set),
  sP: () => (source),
  wk: () => (state),
  zg: () => (mutable_source)
});
/* ESM import */var esm_env__WEBPACK_IMPORTED_MODULE_8__ = __webpack_require__(832);
/* ESM import */var _runtime_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(513);
/* ESM import */var _equality_js__WEBPACK_IMPORTED_MODULE_7__ = __webpack_require__(576);
/* ESM import */var _client_constants__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(924);
/* ESM import */var _errors_js__WEBPACK_IMPORTED_MODULE_10__ = __webpack_require__(626);
/* ESM import */var _flags_index_js__WEBPACK_IMPORTED_MODULE_9__ = __webpack_require__(817);
/* ESM import */var _dev_tracing_js__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(339);
/* ESM import */var _context_js__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(754);
/* ESM import */var _batch_js__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(410);
/* ESM import */var _proxy_js__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(445);
/* ESM import */var _deriveds_js__WEBPACK_IMPORTED_MODULE_6__ = __webpack_require__(462);
/** @import { Derived, Effect, Source, Value } from '#client' */ 










/** @type {Set<any>} */ let inspect_effects = new Set();
/** @type {Map<Source, any>} */ const old_values = new Map();
/**
 * @param {Set<any>} v
 */ function set_inspect_effects(v) {
    inspect_effects = v;
}
let inspect_effects_deferred = false;
function set_inspect_effects_deferred() {
    inspect_effects_deferred = true;
}
/**
 * @template V
 * @param {V} v
 * @param {Error | null} [stack]
 * @returns {Source<V>}
 */ // TODO rename this to `state` throughout the codebase
function source(v, stack) {
    /** @type {Value} */ var signal = {
        f: 0,
        v,
        reactions: null,
        equals: _equality_js__WEBPACK_IMPORTED_MODULE_7__/* .equals */.aI,
        rv: 0,
        wv: 0
    };
    if (esm_env__WEBPACK_IMPORTED_MODULE_8__/* ["default"] */.A && _flags_index_js__WEBPACK_IMPORTED_MODULE_9__/* .tracing_mode_flag */._G) {
        signal.created = stack ?? (0,_dev_tracing_js__WEBPACK_IMPORTED_MODULE_2__/* .get_stack */.sv)('CreatedAt');
        signal.updated = null;
        signal.set_during_effect = false;
        signal.trace = null;
    }
    return signal;
}
/**
 * @template V
 * @param {V} v
 * @param {Error | null} [stack]
 */ /*#__NO_SIDE_EFFECTS__*/ function state(v, stack) {
    const s = source(v, stack);
    (0,_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .push_reaction_value */.tT)(s);
    return s;
}
/**
 * @template V
 * @param {V} initial_value
 * @param {boolean} [immutable]
 * @returns {Source<V>}
 */ /*#__NO_SIDE_EFFECTS__*/ function mutable_source(initial_value) {
    let immutable = arguments.length > 1 && arguments[1] !== void 0 ? arguments[1] : false, trackable = arguments.length > 2 && arguments[2] !== void 0 ? arguments[2] : true;
    const s = source(initial_value);
    if (!immutable) {
        s.equals = _equality_js__WEBPACK_IMPORTED_MODULE_7__/* .safe_equals */.Og;
    }
    // bind the signal to the component context, in case we need to
    // track updates to trigger beforeUpdate/afterUpdate callbacks
    if (_flags_index_js__WEBPACK_IMPORTED_MODULE_9__/* .legacy_mode_flag */.LM && trackable && _context_js__WEBPACK_IMPORTED_MODULE_3__/* .component_context */.UL !== null && _context_js__WEBPACK_IMPORTED_MODULE_3__/* .component_context.l */.UL.l !== null) {
        var _component_context_l;
        ((_component_context_l = _context_js__WEBPACK_IMPORTED_MODULE_3__/* .component_context.l */.UL.l).s ?? (_component_context_l.s = [])).push(s);
    }
    return s;
}
/**
 * @template V
 * @param {Value<V>} source
 * @param {V} value
 */ function mutate(source, value) {
    set(source, untrack(()=>get(source)));
    return value;
}
/**
 * @template V
 * @param {Source<V>} source
 * @param {V} value
 * @param {boolean} [should_proxy]
 * @returns {V}
 */ function set(source, value) {
    let should_proxy = arguments.length > 2 && arguments[2] !== void 0 ? arguments[2] : false;
    if (_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .active_reaction */.hp !== null && // since we are untracking the function inside `$inspect.with` we need to add this check
    // to ensure we error if state is set inside an inspect effect
    (!_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .untracking */.LW || (_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .active_reaction.f */.hp.f & _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .INSPECT_EFFECT */.T1) !== 0) && (0,_context_js__WEBPACK_IMPORTED_MODULE_3__/* .is_runes */.hH)() && (_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .active_reaction.f */.hp.f & (_client_constants__WEBPACK_IMPORTED_MODULE_1__/* .DERIVED */.mj | _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .BLOCK_EFFECT */.kc | _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .ASYNC */.VD | _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .INSPECT_EFFECT */.T1)) !== 0 && !(_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .current_sources */.Bj === null || _runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .current_sources */.Bj === void 0 ? void 0 : _runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .current_sources.includes */.Bj.includes(source))) {
        _errors_js__WEBPACK_IMPORTED_MODULE_10__/* .state_unsafe_mutation */.rZ();
    }
    let new_value = should_proxy ? (0,_proxy_js__WEBPACK_IMPORTED_MODULE_5__/* .proxy */.B)(value) : value;
    if (esm_env__WEBPACK_IMPORTED_MODULE_8__/* ["default"] */.A) {
        (0,_dev_tracing_js__WEBPACK_IMPORTED_MODULE_2__/* .tag_proxy */._e)(new_value, /** @type {string} */ source.label);
    }
    return internal_set(source, new_value);
}
/**
 * @template V
 * @param {Source<V>} source
 * @param {V} value
 * @returns {V}
 */ function internal_set(source, value) {
    if (!source.equals(value)) {
        var old_value = source.v;
        if (_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .is_destroying_effect */.WI) {
            old_values.set(source, value);
        } else {
            old_values.set(source, old_value);
        }
        source.v = value;
        var batch = _batch_js__WEBPACK_IMPORTED_MODULE_4__/* .Batch.ensure */.lP.ensure();
        batch.capture(source, old_value);
        if (esm_env__WEBPACK_IMPORTED_MODULE_8__/* ["default"] */.A) {
            if (_flags_index_js__WEBPACK_IMPORTED_MODULE_9__/* .tracing_mode_flag */._G || _runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .active_effect */.Fg !== null) {
                const error = (0,_dev_tracing_js__WEBPACK_IMPORTED_MODULE_2__/* .get_stack */.sv)('UpdatedAt');
                if (error !== null) {
                    var _source;
                    (_source = source).updated ?? (_source.updated = new Map());
                    let entry = source.updated.get(error.stack);
                    if (!entry) {
                        entry = {
                            error,
                            count: 0
                        };
                        source.updated.set(error.stack, entry);
                    }
                    entry.count++;
                }
            }
            if (_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .active_effect */.Fg !== null) {
                source.set_during_effect = true;
            }
        }
        if ((source.f & _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .DERIVED */.mj) !== 0) {
            // if we are assigning to a dirty derived we set it to clean/maybe dirty but we also eagerly execute it to track the dependencies
            if ((source.f & _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .DIRTY */.jm) !== 0) {
                (0,_deriveds_js__WEBPACK_IMPORTED_MODULE_6__/* .execute_derived */.w6)(/** @type {Derived} */ source);
            }
            (0,_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .set_signal_status */.TC)(source, (source.f & _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .UNOWNED */.L2) === 0 ? _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .CLEAN */.w_ : _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .MAYBE_DIRTY */.ig);
        }
        source.wv = (0,_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .increment_write_version */.Fq)();
        mark_reactions(source, _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .DIRTY */.jm);
        // It's possible that the current reaction might not have up-to-date dependencies
        // whilst it's actively running. So in the case of ensuring it registers the reaction
        // properly for itself, we need to ensure the current effect actually gets
        // scheduled. i.e: `$effect(() => x++)`
        if ((0,_context_js__WEBPACK_IMPORTED_MODULE_3__/* .is_runes */.hH)() && _runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .active_effect */.Fg !== null && (_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .active_effect.f */.Fg.f & _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .CLEAN */.w_) !== 0 && (_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .active_effect.f */.Fg.f & (_client_constants__WEBPACK_IMPORTED_MODULE_1__/* .BRANCH_EFFECT */.Zr | _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .ROOT_EFFECT */.FV)) === 0) {
            if (_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .untracked_writes */.l_ === null) {
                (0,_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .set_untracked_writes */.S0)([
                    source
                ]);
            } else {
                _runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .untracked_writes.push */.l_.push(source);
            }
        }
        if (esm_env__WEBPACK_IMPORTED_MODULE_8__/* ["default"] */.A && inspect_effects.size > 0 && !inspect_effects_deferred) {
            flush_inspect_effects();
        }
    }
    return value;
}
function flush_inspect_effects() {
    inspect_effects_deferred = false;
    const inspects = Array.from(inspect_effects);
    for (const effect of inspects){
        // Mark clean inspect-effects as maybe dirty and then check their dirtiness
        // instead of just updating the effects - this way we avoid overfiring.
        if ((effect.f & _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .CLEAN */.w_) !== 0) {
            (0,_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .set_signal_status */.TC)(effect, _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .MAYBE_DIRTY */.ig);
        }
        if ((0,_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .is_dirty */.Kj)(effect)) {
            (0,_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .update_effect */.gJ)(effect);
        }
    }
    inspect_effects.clear();
}
/**
 * @template {number | bigint} T
 * @param {Source<T>} source
 * @param {1 | -1} [d]
 * @returns {T}
 */ function update(source) {
    let d = arguments.length > 1 && arguments[1] !== void 0 ? arguments[1] : 1;
    var value = get(source);
    var result = d === 1 ? value++ : value--;
    set(source, value);
    // @ts-expect-error
    return result;
}
/**
 * @template {number | bigint} T
 * @param {Source<T>} source
 * @param {1 | -1} [d]
 * @returns {T}
 */ function update_pre(source) {
    let d = arguments.length > 1 && arguments[1] !== void 0 ? arguments[1] : 1;
    var value = get(source);
    // @ts-expect-error
    return set(source, d === 1 ? ++value : --value);
}
/**
 * Silently (without using `get`) increment a source
 * @param {Source<number>} source
 */ function increment(source) {
    set(source, source.v + 1);
}
/**
 * @param {Value} signal
 * @param {number} status should be DIRTY or MAYBE_DIRTY
 * @returns {void}
 */ function mark_reactions(signal, status) {
    var reactions = signal.reactions;
    if (reactions === null) return;
    var runes = (0,_context_js__WEBPACK_IMPORTED_MODULE_3__/* .is_runes */.hH)();
    var length = reactions.length;
    for(var i = 0; i < length; i++){
        var reaction = reactions[i];
        var flags = reaction.f;
        // In legacy mode, skip the current effect to prevent infinite loops
        if (!runes && reaction === _runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .active_effect */.Fg) continue;
        // Inspect effects need to run immediately, so that the stack trace makes sense
        if (esm_env__WEBPACK_IMPORTED_MODULE_8__/* ["default"] */.A && (flags & _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .INSPECT_EFFECT */.T1) !== 0) {
            inspect_effects.add(reaction);
            continue;
        }
        var not_dirty = (flags & _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .DIRTY */.jm) === 0;
        // don't set a DIRTY reaction to MAYBE_DIRTY
        if (not_dirty) {
            (0,_runtime_js__WEBPACK_IMPORTED_MODULE_0__/* .set_signal_status */.TC)(reaction, status);
        }
        if ((flags & _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .DERIVED */.mj) !== 0) {
            mark_reactions(/** @type {Derived} */ reaction, _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .MAYBE_DIRTY */.ig);
        } else if (not_dirty) {
            if ((flags & _client_constants__WEBPACK_IMPORTED_MODULE_1__/* .BLOCK_EFFECT */.kc) !== 0) {
                if (_batch_js__WEBPACK_IMPORTED_MODULE_4__/* .eager_block_effects */.es !== null) {
                    _batch_js__WEBPACK_IMPORTED_MODULE_4__/* .eager_block_effects.push */.es.push(/** @type {Effect} */ reaction);
                }
            }
            (0,_batch_js__WEBPACK_IMPORTED_MODULE_4__/* .schedule_effect */.ec)(/** @type {Effect} */ reaction);
        }
    }
}


}),
485: (function (__unused_webpack_module, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  Or: () => (mount),
  Qv: () => (hydrate),
  j: () => (set_text),
  vs: () => (unmount)
});
/* ESM import */var esm_env__WEBPACK_IMPORTED_MODULE_15__ = __webpack_require__(832);
/* ESM import */var _dom_operations_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(518);
/* ESM import */var _constants_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(178);
/* ESM import */var _runtime_js__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(513);
/* ESM import */var _context_js__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(754);
/* ESM import */var _reactivity_effects_js__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(480);
/* ESM import */var _dom_hydration_js__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(452);
/* ESM import */var _shared_utils_js__WEBPACK_IMPORTED_MODULE_6__ = __webpack_require__(986);
/* ESM import */var _dom_elements_events_js__WEBPACK_IMPORTED_MODULE_7__ = __webpack_require__(417);
/* ESM import */var _dom_blocks_svelte_head_js__WEBPACK_IMPORTED_MODULE_8__ = __webpack_require__(777);
/* ESM import */var _warnings_js__WEBPACK_IMPORTED_MODULE_14__ = __webpack_require__(32);
/* ESM import */var _errors_js__WEBPACK_IMPORTED_MODULE_13__ = __webpack_require__(626);
/* ESM import */var _dom_template_js__WEBPACK_IMPORTED_MODULE_9__ = __webpack_require__(782);
/* ESM import */var _utils_js__WEBPACK_IMPORTED_MODULE_10__ = __webpack_require__(314);
/* ESM import */var _constants_js__WEBPACK_IMPORTED_MODULE_11__ = __webpack_require__(924);
/* ESM import */var _dom_blocks_boundary_js__WEBPACK_IMPORTED_MODULE_12__ = __webpack_require__(899);
/** @import { ComponentContext, Effect, TemplateNode } from '#client' */ /** @import { Component, ComponentType, SvelteComponent, MountOptions } from '../../index.js' */ 















/**
 * This is normally true — block effects should run their intro transitions —
 * but is false during hydration (unless `options.intro` is `true`) and
 * when creating the children of a `<svelte:element>` that just changed tag
 */ let should_intro = true;
/** @param {boolean} value */ function set_should_intro(value) {
    should_intro = value;
}
/**
 * @param {Element} text
 * @param {string} value
 * @returns {void}
 */ function set_text(text, value) {
    var _text;
    // For objects, we apply string coercion (which might make things like $state array references in the template reactive) before diffing
    var str = value == null ? '' : typeof value === 'object' ? value + '' : value;
    // @ts-expect-error
    if (str !== ((_text = text).__t ?? (_text.__t = text.nodeValue))) {
        // @ts-expect-error
        text.__t = str;
        text.nodeValue = str + '';
    }
}
/**
 * Mounts a component to the given target and returns the exports and potentially the props (if compiled with `accessors: true`) of the component.
 * Transitions will play during the initial render unless the `intro` option is set to `false`.
 *
 * @template {Record<string, any>} Props
 * @template {Record<string, any>} Exports
 * @param {ComponentType<SvelteComponent<Props>> | Component<Props, Exports, any>} component
 * @param {MountOptions<Props>} options
 * @returns {Exports}
 */ function mount(component, options) {
    return _mount(component, options);
}
/**
 * Hydrates a component on the given target and returns the exports and potentially the props (if compiled with `accessors: true`) of the component
 *
 * @template {Record<string, any>} Props
 * @template {Record<string, any>} Exports
 * @param {ComponentType<SvelteComponent<Props>> | Component<Props, Exports, any>} component
 * @param {{} extends Props ? {
 * 		target: Document | Element | ShadowRoot;
 * 		props?: Props;
 * 		events?: Record<string, (e: any) => any>;
 *  	context?: Map<any, any>;
 * 		intro?: boolean;
 * 		recover?: boolean;
 * 	} : {
 * 		target: Document | Element | ShadowRoot;
 * 		props: Props;
 * 		events?: Record<string, (e: any) => any>;
 *  	context?: Map<any, any>;
 * 		intro?: boolean;
 * 		recover?: boolean;
 * 	}} options
 * @returns {Exports}
 */ function hydrate(component, options) {
    (0,_dom_operations_js__WEBPACK_IMPORTED_MODULE_0__/* .init_operations */.Ey)();
    options.intro = options.intro ?? false;
    const target = options.target;
    const was_hydrating = _dom_hydration_js__WEBPACK_IMPORTED_MODULE_5__/* .hydrating */.fE;
    const previous_hydrate_node = _dom_hydration_js__WEBPACK_IMPORTED_MODULE_5__/* .hydrate_node */.Xb;
    try {
        var anchor = /** @type {TemplateNode} */ (0,_dom_operations_js__WEBPACK_IMPORTED_MODULE_0__/* .get_first_child */.Zj)(target);
        while(anchor && (anchor.nodeType !== _constants_js__WEBPACK_IMPORTED_MODULE_11__/* .COMMENT_NODE */.dz || /** @type {Comment} */ anchor.data !== _constants_js__WEBPACK_IMPORTED_MODULE_1__/* .HYDRATION_START */.CD)){
            anchor = /** @type {TemplateNode} */ (0,_dom_operations_js__WEBPACK_IMPORTED_MODULE_0__/* .get_next_sibling */.M$)(anchor);
        }
        if (!anchor) {
            throw _constants_js__WEBPACK_IMPORTED_MODULE_1__/* .HYDRATION_ERROR */.kD;
        }
        (0,_dom_hydration_js__WEBPACK_IMPORTED_MODULE_5__/* .set_hydrating */.mK)(true);
        (0,_dom_hydration_js__WEBPACK_IMPORTED_MODULE_5__/* .set_hydrate_node */.W0)(/** @type {Comment} */ anchor);
        const instance = _mount(component, {
            ...options,
            anchor
        });
        (0,_dom_hydration_js__WEBPACK_IMPORTED_MODULE_5__/* .set_hydrating */.mK)(false);
        return /**  @type {Exports} */ instance;
    } catch (error) {
        // re-throw Svelte errors - they are certainly not related to hydration
        if (error instanceof Error && error.message.split('\n').some((line)=>line.startsWith('https://svelte.dev/e/'))) {
            throw error;
        }
        if (error !== _constants_js__WEBPACK_IMPORTED_MODULE_1__/* .HYDRATION_ERROR */.kD) {
            // eslint-disable-next-line no-console
            console.warn('Failed to hydrate: ', error);
        }
        if (options.recover === false) {
            _errors_js__WEBPACK_IMPORTED_MODULE_13__/* .hydration_failed */.Vv();
        }
        // If an error occured above, the operations might not yet have been initialised.
        (0,_dom_operations_js__WEBPACK_IMPORTED_MODULE_0__/* .init_operations */.Ey)();
        (0,_dom_operations_js__WEBPACK_IMPORTED_MODULE_0__/* .clear_text_content */.MC)(target);
        (0,_dom_hydration_js__WEBPACK_IMPORTED_MODULE_5__/* .set_hydrating */.mK)(false);
        return mount(component, options);
    } finally{
        (0,_dom_hydration_js__WEBPACK_IMPORTED_MODULE_5__/* .set_hydrating */.mK)(was_hydrating);
        (0,_dom_hydration_js__WEBPACK_IMPORTED_MODULE_5__/* .set_hydrate_node */.W0)(previous_hydrate_node);
        (0,_dom_blocks_svelte_head_js__WEBPACK_IMPORTED_MODULE_8__/* .reset_head_anchor */.j)();
    }
}
/** @type {Map<string, number>} */ const document_listeners = new Map();
/**
 * @template {Record<string, any>} Exports
 * @param {ComponentType<SvelteComponent<any>> | Component<any>} Component
 * @param {MountOptions} options
 * @returns {Exports}
 */ function _mount(Component, param) {
    let { target, anchor, props = {}, events, context, intro = true } = param;
    (0,_dom_operations_js__WEBPACK_IMPORTED_MODULE_0__/* .init_operations */.Ey)();
    /** @type {Set<string>} */ var registered_events = new Set();
    /** @param {Array<string>} events */ var event_handle = (events)=>{
        for(var i = 0; i < events.length; i++){
            var event_name = events[i];
            if (registered_events.has(event_name)) continue;
            registered_events.add(event_name);
            var passive = (0,_utils_js__WEBPACK_IMPORTED_MODULE_10__/* .is_passive_event */.GY)(event_name);
            // Add the event listener to both the container and the document.
            // The container listener ensures we catch events from within in case
            // the outer content stops propagation of the event.
            target.addEventListener(event_name, _dom_elements_events_js__WEBPACK_IMPORTED_MODULE_7__/* .handle_event_propagation */.n7, {
                passive
            });
            var n = document_listeners.get(event_name);
            if (n === undefined) {
                // The document listener ensures we catch events that originate from elements that were
                // manually moved outside of the container (e.g. via manual portals).
                document.addEventListener(event_name, _dom_elements_events_js__WEBPACK_IMPORTED_MODULE_7__/* .handle_event_propagation */.n7, {
                    passive
                });
                document_listeners.set(event_name, 1);
            } else {
                document_listeners.set(event_name, n + 1);
            }
        }
    };
    event_handle((0,_shared_utils_js__WEBPACK_IMPORTED_MODULE_6__/* .array_from */.bg)(_dom_elements_events_js__WEBPACK_IMPORTED_MODULE_7__/* .all_registered_events */.Ts));
    _dom_elements_events_js__WEBPACK_IMPORTED_MODULE_7__/* .root_event_handles.add */.Sr.add(event_handle);
    /** @type {Exports} */ // @ts-expect-error will be defined because the render effect runs synchronously
    var component = undefined;
    var unmount = (0,_reactivity_effects_js__WEBPACK_IMPORTED_MODULE_4__/* .component_root */.x4)(()=>{
        var anchor_node = anchor ?? target.appendChild((0,_dom_operations_js__WEBPACK_IMPORTED_MODULE_0__/* .create_text */.Pb)());
        (0,_dom_blocks_boundary_js__WEBPACK_IMPORTED_MODULE_12__/* .boundary */.pP)(/** @type {TemplateNode} */ anchor_node, {
            pending: ()=>{}
        }, (anchor_node)=>{
            if (context) {
                (0,_context_js__WEBPACK_IMPORTED_MODULE_3__/* .push */.VC)({});
                var ctx = /** @type {ComponentContext} */ _context_js__WEBPACK_IMPORTED_MODULE_3__/* .component_context */.UL;
                ctx.c = context;
            }
            if (events) {
                // We can't spread the object or else we'd lose the state proxy stuff, if it is one
                /** @type {any} */ props.$$events = events;
            }
            if (_dom_hydration_js__WEBPACK_IMPORTED_MODULE_5__/* .hydrating */.fE) {
                (0,_dom_template_js__WEBPACK_IMPORTED_MODULE_9__/* .assign_nodes */.mX)(/** @type {TemplateNode} */ anchor_node, null);
            }
            should_intro = intro;
            // @ts-expect-error the public typings are not what the actual function looks like
            component = Component(anchor_node, props) || {};
            should_intro = true;
            if (_dom_hydration_js__WEBPACK_IMPORTED_MODULE_5__/* .hydrating */.fE) {
                /** @type {Effect} */ _runtime_js__WEBPACK_IMPORTED_MODULE_2__/* .active_effect.nodes_end */.Fg.nodes_end = _dom_hydration_js__WEBPACK_IMPORTED_MODULE_5__/* .hydrate_node */.Xb;
                if (_dom_hydration_js__WEBPACK_IMPORTED_MODULE_5__/* .hydrate_node */.Xb === null || _dom_hydration_js__WEBPACK_IMPORTED_MODULE_5__/* .hydrate_node.nodeType */.Xb.nodeType !== _constants_js__WEBPACK_IMPORTED_MODULE_11__/* .COMMENT_NODE */.dz || /** @type {Comment} */ _dom_hydration_js__WEBPACK_IMPORTED_MODULE_5__/* .hydrate_node.data */.Xb.data !== _constants_js__WEBPACK_IMPORTED_MODULE_1__/* .HYDRATION_END */.Lc) {
                    _warnings_js__WEBPACK_IMPORTED_MODULE_14__/* .hydration_mismatch */.eZ();
                    throw _constants_js__WEBPACK_IMPORTED_MODULE_1__/* .HYDRATION_ERROR */.kD;
                }
            }
            if (context) {
                (0,_context_js__WEBPACK_IMPORTED_MODULE_3__/* .pop */.uY)();
            }
        });
        return ()=>{
            for (var event_name of registered_events){
                target.removeEventListener(event_name, _dom_elements_events_js__WEBPACK_IMPORTED_MODULE_7__/* .handle_event_propagation */.n7);
                var n = /** @type {number} */ document_listeners.get(event_name);
                if (--n === 0) {
                    document.removeEventListener(event_name, _dom_elements_events_js__WEBPACK_IMPORTED_MODULE_7__/* .handle_event_propagation */.n7);
                    document_listeners.delete(event_name);
                } else {
                    document_listeners.set(event_name, n);
                }
            }
            _dom_elements_events_js__WEBPACK_IMPORTED_MODULE_7__/* .root_event_handles["delete"] */.Sr["delete"](event_handle);
            if (anchor_node !== anchor) {
                var _anchor_node_parentNode;
                (_anchor_node_parentNode = anchor_node.parentNode) === null || _anchor_node_parentNode === void 0 ? void 0 : _anchor_node_parentNode.removeChild(anchor_node);
            }
        };
    });
    mounted_components.set(component, unmount);
    return component;
}
/**
 * References of the components that were mounted or hydrated.
 * Uses a `WeakMap` to avoid memory leaks.
 */ let mounted_components = new WeakMap();
/**
 * Unmounts a component that was previously mounted using `mount` or `hydrate`.
 *
 * Since 5.13.0, if `options.outro` is `true`, [transitions](https://svelte.dev/docs/svelte/transition) will play before the component is removed from the DOM.
 *
 * Returns a `Promise` that resolves after transitions have completed if `options.outro` is true, or immediately otherwise (prior to 5.13.0, returns `void`).
 *
 * ```js
 * import { mount, unmount } from 'svelte';
 * import App from './App.svelte';
 *
 * const app = mount(App, { target: document.body });
 *
 * // later...
 * unmount(app, { outro: true });
 * ```
 * @param {Record<string, any>} component
 * @param {{ outro?: boolean }} [options]
 * @returns {Promise<void>}
 */ function unmount(component, options) {
    const fn = mounted_components.get(component);
    if (fn) {
        mounted_components.delete(component);
        return fn(options);
    }
    if (esm_env__WEBPACK_IMPORTED_MODULE_15__/* ["default"] */.A) {
        if (_constants_js__WEBPACK_IMPORTED_MODULE_11__/* .STATE_SYMBOL */.x3 in component) {
            _warnings_js__WEBPACK_IMPORTED_MODULE_14__/* .state_proxy_unmount */.mx();
        } else {
            _warnings_js__WEBPACK_IMPORTED_MODULE_14__/* .lifecycle_double_unmount */.YY();
        }
    }
    return Promise.resolve();
}


}),
513: (function (__unused_webpack_module, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  BI: () => (set_is_updating_effect),
  Bj: () => (current_sources),
  Fg: () => (active_effect),
  Fq: () => (increment_write_version),
  G0: () => (set_active_reaction),
  Jt: () => (get),
  Kj: () => (is_dirty),
  LW: () => (untracking),
  S0: () => (set_untracked_writes),
  TC: () => (set_signal_status),
  U9: () => (skip_reaction),
  WI: () => (is_destroying_effect),
  eB: () => (set_update_version),
  fT: () => (set_is_destroying_effect),
  gJ: () => (update_effect),
  gU: () => (set_active_effect),
  hp: () => (active_reaction),
  iT: () => (deep_read_state),
  l_: () => (untracked_writes),
  mj: () => (update_reaction),
  pJ: () => (update_version),
  st: () => (is_updating_effect),
  tT: () => (push_reaction_value),
  vz: () => (untrack),
  yR: () => (remove_reactions)
});
/* ESM import */var esm_env__WEBPACK_IMPORTED_MODULE_13__ = __webpack_require__(832);
/* ESM import */var _shared_utils_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(986);
/* ESM import */var _reactivity_effects_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(480);
/* ESM import */var _constants_js__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(924);
/* ESM import */var _reactivity_sources_js__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(264);
/* ESM import */var _reactivity_deriveds_js__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(462);
/* ESM import */var _flags_index_js__WEBPACK_IMPORTED_MODULE_12__ = __webpack_require__(817);
/* ESM import */var _dev_tracing_js__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(339);
/* ESM import */var _context_js__WEBPACK_IMPORTED_MODULE_6__ = __webpack_require__(754);
/* ESM import */var _warnings_js__WEBPACK_IMPORTED_MODULE_14__ = __webpack_require__(32);
/* ESM import */var _reactivity_batch_js__WEBPACK_IMPORTED_MODULE_7__ = __webpack_require__(410);
/* ESM import */var _error_handling_js__WEBPACK_IMPORTED_MODULE_8__ = __webpack_require__(621);
/* ESM import */var _constants_js__WEBPACK_IMPORTED_MODULE_9__ = __webpack_require__(178);
/* ESM import */var _legacy_js__WEBPACK_IMPORTED_MODULE_10__ = __webpack_require__(582);
/* ESM import */var _dom_elements_bindings_shared_js__WEBPACK_IMPORTED_MODULE_11__ = __webpack_require__(408);
/** @import { Derived, Effect, Reaction, Signal, Source, Value } from '#client' */ 














let is_updating_effect = false;
/** @param {boolean} value */ function set_is_updating_effect(value) {
    is_updating_effect = value;
}
let is_destroying_effect = false;
/** @param {boolean} value */ function set_is_destroying_effect(value) {
    is_destroying_effect = value;
}
/** @type {null | Reaction} */ let active_reaction = null;
let untracking = false;
/** @param {null | Reaction} reaction */ function set_active_reaction(reaction) {
    active_reaction = reaction;
}
/** @type {null | Effect} */ let active_effect = null;
/** @param {null | Effect} effect */ function set_active_effect(effect) {
    active_effect = effect;
}
/**
 * When sources are created within a reaction, reading and writing
 * them within that reaction should not cause a re-run
 * @type {null | Source[]}
 */ let current_sources = null;
/** @param {Value} value */ function push_reaction_value(value) {
    if (active_reaction !== null && (!_flags_index_js__WEBPACK_IMPORTED_MODULE_12__/* .async_mode_flag */.I0 || (active_reaction.f & _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .DERIVED */.mj) !== 0)) {
        if (current_sources === null) {
            current_sources = [
                value
            ];
        } else {
            current_sources.push(value);
        }
    }
}
/**
 * The dependencies of the reaction that is currently being executed. In many cases,
 * the dependencies are unchanged between runs, and so this will be `null` unless
 * and until a new dependency is accessed — we track this via `skipped_deps`
 * @type {null | Value[]}
 */ let new_deps = null;
let skipped_deps = 0;
/**
 * Tracks writes that the effect it's executed in doesn't listen to yet,
 * so that the dependency can be added to the effect later on if it then reads it
 * @type {null | Source[]}
 */ let untracked_writes = null;
/** @param {null | Source[]} value */ function set_untracked_writes(value) {
    untracked_writes = value;
}
/**
 * @type {number} Used by sources and deriveds for handling updates.
 * Version starts from 1 so that unowned deriveds differentiate between a created effect and a run one for tracing
 **/ let write_version = 1;
/** @type {number} Used to version each read of a source of derived to avoid duplicating depedencies inside a reaction */ let read_version = 0;
let update_version = read_version;
/** @param {number} value */ function set_update_version(value) {
    update_version = value;
}
// If we are working with a get() chain that has no active container,
// to prevent memory leaks, we skip adding the reaction.
let skip_reaction = false;
function increment_write_version() {
    return ++write_version;
}
/**
 * Determines whether a derived or effect is dirty.
 * If it is MAYBE_DIRTY, will set the status to CLEAN
 * @param {Reaction} reaction
 * @returns {boolean}
 */ function is_dirty(reaction) {
    var flags = reaction.f;
    if ((flags & _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .DIRTY */.jm) !== 0) {
        return true;
    }
    if ((flags & _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .MAYBE_DIRTY */.ig) !== 0) {
        var dependencies = reaction.deps;
        var is_unowned = (flags & _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .UNOWNED */.L2) !== 0;
        if (dependencies !== null) {
            var i;
            var dependency;
            var is_disconnected = (flags & _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .DISCONNECTED */._N) !== 0;
            var is_unowned_connected = is_unowned && active_effect !== null && !skip_reaction;
            var length = dependencies.length;
            // If we are working with a disconnected or an unowned signal that is now connected (due to an active effect)
            // then we need to re-connect the reaction to the dependency, unless the effect has already been destroyed
            // (which can happen if the derived is read by an async derived)
            if ((is_disconnected || is_unowned_connected) && (active_effect === null || (active_effect.f & _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .DESTROYED */.o5) === 0)) {
                var derived = /** @type {Derived} */ reaction;
                var parent = derived.parent;
                for(i = 0; i < length; i++){
                    var _dependency_reactions;
                    dependency = dependencies[i];
                    // We always re-add all reactions (even duplicates) if the derived was
                    // previously disconnected, however we don't if it was unowned as we
                    // de-duplicate dependencies in that case
                    if (is_disconnected || !(dependency === null || dependency === void 0 ? void 0 : (_dependency_reactions = dependency.reactions) === null || _dependency_reactions === void 0 ? void 0 : _dependency_reactions.includes(derived))) {
                        var _dependency;
                        ((_dependency = dependency).reactions ?? (_dependency.reactions = [])).push(derived);
                    }
                }
                if (is_disconnected) {
                    derived.f ^= _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .DISCONNECTED */._N;
                }
                // If the unowned derived is now fully connected to the graph again (it's unowned and reconnected, has a parent
                // and the parent is not unowned), then we can mark it as connected again, removing the need for the unowned
                // flag
                if (is_unowned_connected && parent !== null && (parent.f & _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .UNOWNED */.L2) === 0) {
                    derived.f ^= _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .UNOWNED */.L2;
                }
            }
            for(i = 0; i < length; i++){
                dependency = dependencies[i];
                if (is_dirty(/** @type {Derived} */ dependency)) {
                    (0,_reactivity_deriveds_js__WEBPACK_IMPORTED_MODULE_4__/* .update_derived */.c2)(/** @type {Derived} */ dependency);
                }
                if (dependency.wv > reaction.wv) {
                    return true;
                }
            }
        }
        // Unowned signals should never be marked as clean unless they
        // are used within an active_effect without skip_reaction
        if (!is_unowned || active_effect !== null && !skip_reaction) {
            set_signal_status(reaction, _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .CLEAN */.w_);
        }
    }
    return false;
}
/**
 * @param {Value} signal
 * @param {Effect} effect
 * @param {boolean} [root]
 */ function schedule_possible_effect_self_invalidation(signal, effect) {
    let root = arguments.length > 2 && arguments[2] !== void 0 ? arguments[2] : true;
    var reactions = signal.reactions;
    if (reactions === null) return;
    if (!_flags_index_js__WEBPACK_IMPORTED_MODULE_12__/* .async_mode_flag */.I0 && (current_sources === null || current_sources === void 0 ? void 0 : current_sources.includes(signal))) {
        return;
    }
    for(var i = 0; i < reactions.length; i++){
        var reaction = reactions[i];
        if ((reaction.f & _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .DERIVED */.mj) !== 0) {
            schedule_possible_effect_self_invalidation(/** @type {Derived} */ reaction, effect, false);
        } else if (effect === reaction) {
            if (root) {
                set_signal_status(reaction, _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .DIRTY */.jm);
            } else if ((reaction.f & _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .CLEAN */.w_) !== 0) {
                set_signal_status(reaction, _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .MAYBE_DIRTY */.ig);
            }
            (0,_reactivity_batch_js__WEBPACK_IMPORTED_MODULE_7__/* .schedule_effect */.ec)(/** @type {Effect} */ reaction);
        }
    }
}
/** @param {Reaction} reaction */ function update_reaction(reaction) {
    var previous_deps = new_deps;
    var previous_skipped_deps = skipped_deps;
    var previous_untracked_writes = untracked_writes;
    var previous_reaction = active_reaction;
    var previous_skip_reaction = skip_reaction;
    var previous_sources = current_sources;
    var previous_component_context = _context_js__WEBPACK_IMPORTED_MODULE_6__/* .component_context */.UL;
    var previous_untracking = untracking;
    var previous_update_version = update_version;
    var flags = reaction.f;
    new_deps = /** @type {null | Value[]} */ null;
    skipped_deps = 0;
    untracked_writes = null;
    skip_reaction = (flags & _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .UNOWNED */.L2) !== 0 && (untracking || !is_updating_effect || active_reaction === null);
    active_reaction = (flags & (_constants_js__WEBPACK_IMPORTED_MODULE_2__/* .BRANCH_EFFECT */.Zr | _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .ROOT_EFFECT */.FV)) === 0 ? reaction : null;
    current_sources = null;
    (0,_context_js__WEBPACK_IMPORTED_MODULE_6__/* .set_component_context */.De)(reaction.ctx);
    untracking = false;
    update_version = ++read_version;
    if (reaction.ac !== null) {
        (0,_dom_elements_bindings_shared_js__WEBPACK_IMPORTED_MODULE_11__/* .without_reactive_context */.$w)(()=>{
            /** @type {AbortController} */ reaction.ac.abort(_constants_js__WEBPACK_IMPORTED_MODULE_2__/* .STALE_REACTION */.In);
        });
        reaction.ac = null;
    }
    try {
        reaction.f |= _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .REACTION_IS_UPDATING */.EY;
        var fn = /** @type {Function} */ reaction.fn;
        var result = fn();
        var deps = reaction.deps;
        if (new_deps !== null) {
            var i;
            remove_reactions(reaction, skipped_deps);
            if (deps !== null && skipped_deps > 0) {
                deps.length = skipped_deps + new_deps.length;
                for(i = 0; i < new_deps.length; i++){
                    deps[skipped_deps + i] = new_deps[i];
                }
            } else {
                reaction.deps = deps = new_deps;
            }
            if (!skip_reaction || // Deriveds that already have reactions can cleanup, so we still add them as reactions
            (flags & _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .DERIVED */.mj) !== 0 && /** @type {import('#client').Derived} */ reaction.reactions !== null) {
                for(i = skipped_deps; i < deps.length; i++){
                    var _deps_i;
                    ((_deps_i = deps[i]).reactions ?? (_deps_i.reactions = [])).push(reaction);
                }
            }
        } else if (deps !== null && skipped_deps < deps.length) {
            remove_reactions(reaction, skipped_deps);
            deps.length = skipped_deps;
        }
        // If we're inside an effect and we have untracked writes, then we need to
        // ensure that if any of those untracked writes result in re-invalidation
        // of the current effect, then that happens accordingly
        if ((0,_context_js__WEBPACK_IMPORTED_MODULE_6__/* .is_runes */.hH)() && untracked_writes !== null && !untracking && deps !== null && (reaction.f & (_constants_js__WEBPACK_IMPORTED_MODULE_2__/* .DERIVED */.mj | _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .MAYBE_DIRTY */.ig | _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .DIRTY */.jm)) === 0) {
            for(i = 0; i < /** @type {Source[]} */ untracked_writes.length; i++){
                schedule_possible_effect_self_invalidation(untracked_writes[i], /** @type {Effect} */ reaction);
            }
        }
        // If we are returning to an previous reaction then
        // we need to increment the read version to ensure that
        // any dependencies in this reaction aren't marked with
        // the same version
        if (previous_reaction !== null && previous_reaction !== reaction) {
            read_version++;
            if (untracked_writes !== null) {
                if (previous_untracked_writes === null) {
                    previous_untracked_writes = untracked_writes;
                } else {
                    previous_untracked_writes.push(.../** @type {Source[]} */ untracked_writes);
                }
            }
        }
        if ((reaction.f & _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .ERROR_VALUE */.dH) !== 0) {
            reaction.f ^= _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .ERROR_VALUE */.dH;
        }
        return result;
    } catch (error) {
        return (0,_error_handling_js__WEBPACK_IMPORTED_MODULE_8__/* .handle_error */.i)(error);
    } finally{
        reaction.f ^= _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .REACTION_IS_UPDATING */.EY;
        new_deps = previous_deps;
        skipped_deps = previous_skipped_deps;
        untracked_writes = previous_untracked_writes;
        active_reaction = previous_reaction;
        skip_reaction = previous_skip_reaction;
        current_sources = previous_sources;
        (0,_context_js__WEBPACK_IMPORTED_MODULE_6__/* .set_component_context */.De)(previous_component_context);
        untracking = previous_untracking;
        update_version = previous_update_version;
    }
}
/**
 * @template V
 * @param {Reaction} signal
 * @param {Value<V>} dependency
 * @returns {void}
 */ function remove_reaction(signal, dependency) {
    let reactions = dependency.reactions;
    if (reactions !== null) {
        var index = _shared_utils_js__WEBPACK_IMPORTED_MODULE_0__/* .index_of.call */.lc.call(reactions, signal);
        if (index !== -1) {
            var new_length = reactions.length - 1;
            if (new_length === 0) {
                reactions = dependency.reactions = null;
            } else {
                // Swap with last element and then remove.
                reactions[index] = reactions[new_length];
                reactions.pop();
            }
        }
    }
    // If the derived has no reactions, then we can disconnect it from the graph,
    // allowing it to either reconnect in the future, or be GC'd by the VM.
    if (reactions === null && (dependency.f & _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .DERIVED */.mj) !== 0 && // Destroying a child effect while updating a parent effect can cause a dependency to appear
    // to be unused, when in fact it is used by the currently-updating parent. Checking `new_deps`
    // allows us to skip the expensive work of disconnecting and immediately reconnecting it
    (new_deps === null || !new_deps.includes(dependency))) {
        set_signal_status(dependency, _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .MAYBE_DIRTY */.ig);
        // If we are working with a derived that is owned by an effect, then mark it as being
        // disconnected.
        if ((dependency.f & (_constants_js__WEBPACK_IMPORTED_MODULE_2__/* .UNOWNED */.L2 | _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .DISCONNECTED */._N)) === 0) {
            dependency.f ^= _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .DISCONNECTED */._N;
        }
        // Disconnect any reactions owned by this reaction
        (0,_reactivity_deriveds_js__WEBPACK_IMPORTED_MODULE_4__/* .destroy_derived_effects */.ge)(/** @type {Derived} **/ dependency);
        remove_reactions(/** @type {Derived} **/ dependency, 0);
    }
}
/**
 * @param {Reaction} signal
 * @param {number} start_index
 * @returns {void}
 */ function remove_reactions(signal, start_index) {
    var dependencies = signal.deps;
    if (dependencies === null) return;
    for(var i = start_index; i < dependencies.length; i++){
        remove_reaction(signal, dependencies[i]);
    }
}
/**
 * @param {Effect} effect
 * @returns {void}
 */ function update_effect(effect) {
    var flags = effect.f;
    if ((flags & _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .DESTROYED */.o5) !== 0) {
        return;
    }
    set_signal_status(effect, _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .CLEAN */.w_);
    var previous_effect = active_effect;
    var was_updating_effect = is_updating_effect;
    active_effect = effect;
    is_updating_effect = true;
    if (esm_env__WEBPACK_IMPORTED_MODULE_13__/* ["default"] */.A) {
        var previous_component_fn = _context_js__WEBPACK_IMPORTED_MODULE_6__/* .dev_current_component_function */.DE;
        (0,_context_js__WEBPACK_IMPORTED_MODULE_6__/* .set_dev_current_component_function */.Mo)(effect.component_function);
        var previous_stack = /** @type {any} */ _context_js__WEBPACK_IMPORTED_MODULE_6__/* .dev_stack */.lv;
        // only block effects have a dev stack, keep the current one otherwise
        (0,_context_js__WEBPACK_IMPORTED_MODULE_6__/* .set_dev_stack */.O2)(effect.dev_stack ?? _context_js__WEBPACK_IMPORTED_MODULE_6__/* .dev_stack */.lv);
    }
    try {
        if ((flags & _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .BLOCK_EFFECT */.kc) !== 0) {
            (0,_reactivity_effects_js__WEBPACK_IMPORTED_MODULE_1__/* .destroy_block_effect_children */.pk)(effect);
        } else {
            (0,_reactivity_effects_js__WEBPACK_IMPORTED_MODULE_1__/* .destroy_effect_children */.F3)(effect);
        }
        (0,_reactivity_effects_js__WEBPACK_IMPORTED_MODULE_1__/* .execute_effect_teardown */.Nq)(effect);
        var teardown = update_reaction(effect);
        effect.teardown = typeof teardown === 'function' ? teardown : null;
        effect.wv = write_version;
        // In DEV, increment versions of any sources that were written to during the effect,
        // so that they are correctly marked as dirty when the effect re-runs
        if (esm_env__WEBPACK_IMPORTED_MODULE_13__/* ["default"] */.A && _flags_index_js__WEBPACK_IMPORTED_MODULE_12__/* .tracing_mode_flag */._G && (effect.f & _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .DIRTY */.jm) !== 0 && effect.deps !== null) {
            for (var dep of effect.deps){
                if (dep.set_during_effect) {
                    dep.wv = increment_write_version();
                    dep.set_during_effect = false;
                }
            }
        }
    } finally{
        is_updating_effect = was_updating_effect;
        active_effect = previous_effect;
        if (esm_env__WEBPACK_IMPORTED_MODULE_13__/* ["default"] */.A) {
            (0,_context_js__WEBPACK_IMPORTED_MODULE_6__/* .set_dev_current_component_function */.Mo)(previous_component_fn);
            (0,_context_js__WEBPACK_IMPORTED_MODULE_6__/* .set_dev_stack */.O2)(previous_stack);
        }
    }
}
/**
 * Returns a promise that resolves once any pending state changes have been applied.
 * @returns {Promise<void>}
 */ async function tick() {
    if (async_mode_flag) {
        return new Promise((f)=>requestAnimationFrame(()=>f()));
    }
    await Promise.resolve();
    // By calling flushSync we guarantee that any pending state changes are applied after one tick.
    // TODO look into whether we can make flushing subsequent updates synchronously in the future.
    flushSync();
}
/**
 * Returns a promise that resolves once any state changes, and asynchronous work resulting from them,
 * have resolved and the DOM has been updated
 * @returns {Promise<void>}
 * @since 5.36
 */ function settled() {
    return Batch.ensure().settled();
}
/**
 * @template V
 * @param {Value<V>} signal
 * @returns {V}
 */ function get(signal) {
    var flags = signal.f;
    var is_derived = (flags & _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .DERIVED */.mj) !== 0;
    _legacy_js__WEBPACK_IMPORTED_MODULE_10__/* .captured_signals */.J === null || _legacy_js__WEBPACK_IMPORTED_MODULE_10__/* .captured_signals */.J === void 0 ? void 0 : _legacy_js__WEBPACK_IMPORTED_MODULE_10__/* .captured_signals.add */.J.add(signal);
    // Register the dependency on the current reaction signal.
    if (active_reaction !== null && !untracking) {
        // if we're in a derived that is being read inside an _async_ derived,
        // it's possible that the effect was already destroyed. In this case,
        // we don't add the dependency, because that would create a memory leak
        var destroyed = active_effect !== null && (active_effect.f & _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .DESTROYED */.o5) !== 0;
        if (!destroyed && !(current_sources === null || current_sources === void 0 ? void 0 : current_sources.includes(signal))) {
            var deps = active_reaction.deps;
            if ((active_reaction.f & _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .REACTION_IS_UPDATING */.EY) !== 0) {
                // we're in the effect init/update cycle
                if (signal.rv < read_version) {
                    signal.rv = read_version;
                    // If the signal is accessing the same dependencies in the same
                    // order as it did last time, increment `skipped_deps`
                    // rather than updating `new_deps`, which creates GC cost
                    if (new_deps === null && deps !== null && deps[skipped_deps] === signal) {
                        skipped_deps++;
                    } else if (new_deps === null) {
                        new_deps = [
                            signal
                        ];
                    } else if (!skip_reaction || !new_deps.includes(signal)) {
                        // Normally we can push duplicated dependencies to `new_deps`, but if we're inside
                        // an unowned derived because skip_reaction is true, then we need to ensure that
                        // we don't have duplicates
                        new_deps.push(signal);
                    }
                }
            } else {
                var // we're adding a dependency outside the init/update cycle
                // (i.e. after an `await`)
                _active_reaction;
                ((_active_reaction = active_reaction).deps ?? (_active_reaction.deps = [])).push(signal);
                var reactions = signal.reactions;
                if (reactions === null) {
                    signal.reactions = [
                        active_reaction
                    ];
                } else if (!reactions.includes(active_reaction)) {
                    reactions.push(active_reaction);
                }
            }
        }
    } else if (is_derived && /** @type {Derived} */ signal.deps === null && /** @type {Derived} */ signal.effects === null) {
        var derived = /** @type {Derived} */ signal;
        var parent = derived.parent;
        if (parent !== null && (parent.f & _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .UNOWNED */.L2) === 0) {
            // If the derived is owned by another derived then mark it as unowned
            // as the derived value might have been referenced in a different context
            // since and thus its parent might not be its true owner anymore
            derived.f ^= _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .UNOWNED */.L2;
        }
    }
    if (esm_env__WEBPACK_IMPORTED_MODULE_13__/* ["default"] */.A) {
        if (_reactivity_deriveds_js__WEBPACK_IMPORTED_MODULE_4__/* .current_async_effect */.vO) {
            var _current_async_effect_deps;
            var tracking = (_reactivity_deriveds_js__WEBPACK_IMPORTED_MODULE_4__/* .current_async_effect.f */.vO.f & _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .REACTION_IS_UPDATING */.EY) !== 0;
            var was_read = (_current_async_effect_deps = _reactivity_deriveds_js__WEBPACK_IMPORTED_MODULE_4__/* .current_async_effect.deps */.vO.deps) === null || _current_async_effect_deps === void 0 ? void 0 : _current_async_effect_deps.includes(signal);
            if (!tracking && !untracking && !was_read) {
                _warnings_js__WEBPACK_IMPORTED_MODULE_14__/* .await_reactivity_loss */._2(/** @type {string} */ signal.label);
                var trace = (0,_dev_tracing_js__WEBPACK_IMPORTED_MODULE_5__/* .get_stack */.sv)('TracedAt');
                // eslint-disable-next-line no-console
                if (trace) console.warn(trace);
            }
        }
        _reactivity_deriveds_js__WEBPACK_IMPORTED_MODULE_4__/* .recent_async_deriveds["delete"] */.kX["delete"](signal);
        if (_flags_index_js__WEBPACK_IMPORTED_MODULE_12__/* .tracing_mode_flag */._G && !untracking && _dev_tracing_js__WEBPACK_IMPORTED_MODULE_5__/* .tracing_expressions */.ho !== null && active_reaction !== null && _dev_tracing_js__WEBPACK_IMPORTED_MODULE_5__/* .tracing_expressions.reaction */.ho.reaction === active_reaction) {
            // Used when mapping state between special blocks like `each`
            if (signal.trace) {
                signal.trace();
            } else {
                trace = (0,_dev_tracing_js__WEBPACK_IMPORTED_MODULE_5__/* .get_stack */.sv)('TracedAt');
                if (trace) {
                    var entry = _dev_tracing_js__WEBPACK_IMPORTED_MODULE_5__/* .tracing_expressions.entries.get */.ho.entries.get(signal);
                    if (entry === undefined) {
                        entry = {
                            traces: []
                        };
                        _dev_tracing_js__WEBPACK_IMPORTED_MODULE_5__/* .tracing_expressions.entries.set */.ho.entries.set(signal, entry);
                    }
                    var last = entry.traces[entry.traces.length - 1];
                    // traces can be duplicated, e.g. by `snapshot` invoking both
                    // both `getOwnPropertyDescriptor` and `get` traps at once
                    if (trace.stack !== (last === null || last === void 0 ? void 0 : last.stack)) {
                        entry.traces.push(trace);
                    }
                }
            }
        }
    }
    if (is_destroying_effect) {
        if (_reactivity_sources_js__WEBPACK_IMPORTED_MODULE_3__/* .old_values.has */.bJ.has(signal)) {
            return _reactivity_sources_js__WEBPACK_IMPORTED_MODULE_3__/* .old_values.get */.bJ.get(signal);
        }
        if (is_derived) {
            derived = /** @type {Derived} */ signal;
            var value = derived.v;
            // if the derived is dirty and has reactions, or depends on the values that just changed, re-execute
            // (a derived can be maybe_dirty due to the effect destroy removing its last reaction)
            if ((derived.f & _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .CLEAN */.w_) === 0 && derived.reactions !== null || depends_on_old_values(derived)) {
                value = (0,_reactivity_deriveds_js__WEBPACK_IMPORTED_MODULE_4__/* .execute_derived */.w6)(derived);
            }
            _reactivity_sources_js__WEBPACK_IMPORTED_MODULE_3__/* .old_values.set */.bJ.set(derived, value);
            return value;
        }
    } else if (is_derived) {
        derived = /** @type {Derived} */ signal;
        if (_reactivity_batch_js__WEBPACK_IMPORTED_MODULE_7__/* .batch_deriveds */.G1 === null || _reactivity_batch_js__WEBPACK_IMPORTED_MODULE_7__/* .batch_deriveds */.G1 === void 0 ? void 0 : _reactivity_batch_js__WEBPACK_IMPORTED_MODULE_7__/* .batch_deriveds.has */.G1.has(derived)) {
            return _reactivity_batch_js__WEBPACK_IMPORTED_MODULE_7__/* .batch_deriveds.get */.G1.get(derived);
        }
        if (is_dirty(derived)) {
            (0,_reactivity_deriveds_js__WEBPACK_IMPORTED_MODULE_4__/* .update_derived */.c2)(derived);
        }
    }
    if ((signal.f & _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .ERROR_VALUE */.dH) !== 0) {
        throw signal.v;
    }
    return signal.v;
}
/** @param {Derived} derived */ function depends_on_old_values(derived) {
    if (derived.v === _constants_js__WEBPACK_IMPORTED_MODULE_9__/* .UNINITIALIZED */.UP) return true; // we don't know, so assume the worst
    if (derived.deps === null) return false;
    for (const dep of derived.deps){
        if (_reactivity_sources_js__WEBPACK_IMPORTED_MODULE_3__/* .old_values.has */.bJ.has(dep)) {
            return true;
        }
        if ((dep.f & _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .DERIVED */.mj) !== 0 && depends_on_old_values(/** @type {Derived} */ dep)) {
            return true;
        }
    }
    return false;
}
/**
 * Like `get`, but checks for `undefined`. Used for `var` declarations because they can be accessed before being declared
 * @template V
 * @param {Value<V> | undefined} signal
 * @returns {V | undefined}
 */ function safe_get(signal) {
    return signal && get(signal);
}
/**
 * When used inside a [`$derived`](https://svelte.dev/docs/svelte/$derived) or [`$effect`](https://svelte.dev/docs/svelte/$effect),
 * any state read inside `fn` will not be treated as a dependency.
 *
 * ```ts
 * $effect(() => {
 *   // this will run when `data` changes, but not when `time` changes
 *   save(data, {
 *     timestamp: untrack(() => time)
 *   });
 * });
 * ```
 * @template T
 * @param {() => T} fn
 * @returns {T}
 */ function untrack(fn) {
    var previous_untracking = untracking;
    try {
        untracking = true;
        return fn();
    } finally{
        untracking = previous_untracking;
    }
}
const STATUS_MASK = ~(_constants_js__WEBPACK_IMPORTED_MODULE_2__/* .DIRTY */.jm | _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .MAYBE_DIRTY */.ig | _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .CLEAN */.w_);
/**
 * @param {Signal} signal
 * @param {number} status
 * @returns {void}
 */ function set_signal_status(signal, status) {
    signal.f = signal.f & STATUS_MASK | status;
}
/**
 * @param {Record<string, unknown>} obj
 * @param {string[]} keys
 * @returns {Record<string, unknown>}
 */ function exclude_from_object(obj, keys) {
    /** @type {Record<string, unknown>} */ var result = {};
    for(var key in obj){
        if (!keys.includes(key)) {
            result[key] = obj[key];
        }
    }
    return result;
}
/**
 * Possibly traverse an object and read all its properties so that they're all reactive in case this is `$state`.
 * Does only check first level of an object for performance reasons (heuristic should be good for 99% of all cases).
 * @param {any} value
 * @returns {void}
 */ function deep_read_state(value) {
    if (typeof value !== 'object' || !value || value instanceof EventTarget) {
        return;
    }
    if (_constants_js__WEBPACK_IMPORTED_MODULE_2__/* .STATE_SYMBOL */.x3 in value) {
        deep_read(value);
    } else if (!Array.isArray(value)) {
        for(let key in value){
            const prop = value[key];
            if (typeof prop === 'object' && prop && _constants_js__WEBPACK_IMPORTED_MODULE_2__/* .STATE_SYMBOL */.x3 in prop) {
                deep_read(prop);
            }
        }
    }
}
/**
 * Deeply traverse an object and read all its properties
 * so that they're all reactive in case this is `$state`
 * @param {any} value
 * @param {Set<any>} visited
 * @returns {void}
 */ function deep_read(value) {
    let visited = arguments.length > 1 && arguments[1] !== void 0 ? arguments[1] : new Set();
    if (typeof value === 'object' && value !== null && // We don't want to traverse DOM elements
    !(value instanceof EventTarget) && !visited.has(value)) {
        visited.add(value);
        // When working with a possible SvelteDate, this
        // will ensure we capture changes to it.
        if (value instanceof Date) {
            value.getTime();
        }
        for(let key in value){
            try {
                deep_read(value[key], visited);
            } catch (e) {
            // continue
            }
        }
        const proto = (0,_shared_utils_js__WEBPACK_IMPORTED_MODULE_0__/* .get_prototype_of */.Oh)(value);
        if (proto !== Object.prototype && proto !== Array.prototype && proto !== Map.prototype && proto !== Set.prototype && proto !== Date.prototype) {
            const descriptors = (0,_shared_utils_js__WEBPACK_IMPORTED_MODULE_0__/* .get_descriptors */.CL)(proto);
            for(let key in descriptors){
                const get = descriptors[key].get;
                if (get) {
                    try {
                        get.call(value);
                    } catch (e) {
                    // continue
                    }
                }
            }
        }
    }
}


}),
32: (function (__unused_webpack_module, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  CF: () => (svelte_boundary_reset_noop),
  Cy: () => (await_waterfall),
  Y9: () => (hydration_html_changed),
  YY: () => (lifecycle_double_unmount),
  _2: () => (await_reactivity_loss),
  eZ: () => (hydration_mismatch),
  mx: () => (state_proxy_unmount),
  ns: () => (state_proxy_equality_mismatch)
});
/* ESM import */var esm_env__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(832);
/* This file is generated by scripts/process-messages/index.js. Do not edit! */ 
var bold = 'font-weight: bold';
var normal = 'font-weight: normal';
/**
 * Assignment to `%property%` property (%location%) will evaluate to the right-hand side, not the value of `%property%` following the assignment. This may result in unexpected behaviour.
 * @param {string} property
 * @param {string} location
 */ function assignment_value_stale(property, location) {
    if (DEV) {
        console.warn(`%c[svelte] assignment_value_stale\n%cAssignment to \`${property}\` property (${location}) will evaluate to the right-hand side, not the value of \`${property}\` following the assignment. This may result in unexpected behaviour.\nhttps://svelte.dev/e/assignment_value_stale`, bold, normal);
    } else {
        console.warn(`https://svelte.dev/e/assignment_value_stale`);
    }
}
/**
 * Detected reactivity loss when reading `%name%`. This happens when state is read in an async function after an earlier `await`
 * @param {string} name
 */ function await_reactivity_loss(name) {
    if (esm_env__WEBPACK_IMPORTED_MODULE_0__/* ["default"] */.A) {
        console.warn(`%c[svelte] await_reactivity_loss\n%cDetected reactivity loss when reading \`${name}\`. This happens when state is read in an async function after an earlier \`await\`\nhttps://svelte.dev/e/await_reactivity_loss`, bold, normal);
    } else {
        console.warn(`https://svelte.dev/e/await_reactivity_loss`);
    }
}
/**
 * An async derived, `%name%` (%location%) was not read immediately after it resolved. This often indicates an unnecessary waterfall, which can slow down your app
 * @param {string} name
 * @param {string} location
 */ function await_waterfall(name, location) {
    if (esm_env__WEBPACK_IMPORTED_MODULE_0__/* ["default"] */.A) {
        console.warn(`%c[svelte] await_waterfall\n%cAn async derived, \`${name}\` (${location}) was not read immediately after it resolved. This often indicates an unnecessary waterfall, which can slow down your app\nhttps://svelte.dev/e/await_waterfall`, bold, normal);
    } else {
        console.warn(`https://svelte.dev/e/await_waterfall`);
    }
}
/**
 * `%binding%` (%location%) is binding to a non-reactive property
 * @param {string} binding
 * @param {string | undefined | null} [location]
 */ function binding_property_non_reactive(binding, location) {
    if (DEV) {
        console.warn(`%c[svelte] binding_property_non_reactive\n%c${location ? `\`${binding}\` (${location}) is binding to a non-reactive property` : `\`${binding}\` is binding to a non-reactive property`}\nhttps://svelte.dev/e/binding_property_non_reactive`, bold, normal);
    } else {
        console.warn(`https://svelte.dev/e/binding_property_non_reactive`);
    }
}
/**
 * Your `console.%method%` contained `$state` proxies. Consider using `$inspect(...)` or `$state.snapshot(...)` instead
 * @param {string} method
 */ function console_log_state(method) {
    if (DEV) {
        console.warn(`%c[svelte] console_log_state\n%cYour \`console.${method}\` contained \`$state\` proxies. Consider using \`$inspect(...)\` or \`$state.snapshot(...)\` instead\nhttps://svelte.dev/e/console_log_state`, bold, normal);
    } else {
        console.warn(`https://svelte.dev/e/console_log_state`);
    }
}
/**
 * %handler% should be a function. Did you mean to %suggestion%?
 * @param {string} handler
 * @param {string} suggestion
 */ function event_handler_invalid(handler, suggestion) {
    if (DEV) {
        console.warn(`%c[svelte] event_handler_invalid\n%c${handler} should be a function. Did you mean to ${suggestion}?\nhttps://svelte.dev/e/event_handler_invalid`, bold, normal);
    } else {
        console.warn(`https://svelte.dev/e/event_handler_invalid`);
    }
}
/**
 * The `%attribute%` attribute on `%html%` changed its value between server and client renders. The client value, `%value%`, will be ignored in favour of the server value
 * @param {string} attribute
 * @param {string} html
 * @param {string} value
 */ function hydration_attribute_changed(attribute, html, value) {
    if (DEV) {
        console.warn(`%c[svelte] hydration_attribute_changed\n%cThe \`${attribute}\` attribute on \`${html}\` changed its value between server and client renders. The client value, \`${value}\`, will be ignored in favour of the server value\nhttps://svelte.dev/e/hydration_attribute_changed`, bold, normal);
    } else {
        console.warn(`https://svelte.dev/e/hydration_attribute_changed`);
    }
}
/**
 * The value of an `{@html ...}` block %location% changed between server and client renders. The client value will be ignored in favour of the server value
 * @param {string | undefined | null} [location]
 */ function hydration_html_changed(location) {
    if (esm_env__WEBPACK_IMPORTED_MODULE_0__/* ["default"] */.A) {
        console.warn(`%c[svelte] hydration_html_changed\n%c${location ? `The value of an \`{@html ...}\` block ${location} changed between server and client renders. The client value will be ignored in favour of the server value` : 'The value of an `{@html ...}` block changed between server and client renders. The client value will be ignored in favour of the server value'}\nhttps://svelte.dev/e/hydration_html_changed`, bold, normal);
    } else {
        console.warn(`https://svelte.dev/e/hydration_html_changed`);
    }
}
/**
 * Hydration failed because the initial UI does not match what was rendered on the server. The error occurred near %location%
 * @param {string | undefined | null} [location]
 */ function hydration_mismatch(location) {
    if (esm_env__WEBPACK_IMPORTED_MODULE_0__/* ["default"] */.A) {
        console.warn(`%c[svelte] hydration_mismatch\n%c${location ? `Hydration failed because the initial UI does not match what was rendered on the server. The error occurred near ${location}` : 'Hydration failed because the initial UI does not match what was rendered on the server'}\nhttps://svelte.dev/e/hydration_mismatch`, bold, normal);
    } else {
        console.warn(`https://svelte.dev/e/hydration_mismatch`);
    }
}
/**
 * The `render` function passed to `createRawSnippet` should return HTML for a single element
 */ function invalid_raw_snippet_render() {
    if (DEV) {
        console.warn(`%c[svelte] invalid_raw_snippet_render\n%cThe \`render\` function passed to \`createRawSnippet\` should return HTML for a single element\nhttps://svelte.dev/e/invalid_raw_snippet_render`, bold, normal);
    } else {
        console.warn(`https://svelte.dev/e/invalid_raw_snippet_render`);
    }
}
/**
 * Detected a migrated `$:` reactive block in `%filename%` that both accesses and updates the same reactive value. This may cause recursive updates when converted to an `$effect`.
 * @param {string} filename
 */ function legacy_recursive_reactive_block(filename) {
    if (DEV) {
        console.warn(`%c[svelte] legacy_recursive_reactive_block\n%cDetected a migrated \`$:\` reactive block in \`${filename}\` that both accesses and updates the same reactive value. This may cause recursive updates when converted to an \`$effect\`.\nhttps://svelte.dev/e/legacy_recursive_reactive_block`, bold, normal);
    } else {
        console.warn(`https://svelte.dev/e/legacy_recursive_reactive_block`);
    }
}
/**
 * Tried to unmount a component that was not mounted
 */ function lifecycle_double_unmount() {
    if (esm_env__WEBPACK_IMPORTED_MODULE_0__/* ["default"] */.A) {
        console.warn(`%c[svelte] lifecycle_double_unmount\n%cTried to unmount a component that was not mounted\nhttps://svelte.dev/e/lifecycle_double_unmount`, bold, normal);
    } else {
        console.warn(`https://svelte.dev/e/lifecycle_double_unmount`);
    }
}
/**
 * %parent% passed property `%prop%` to %child% with `bind:`, but its parent component %owner% did not declare `%prop%` as a binding. Consider creating a binding between %owner% and %parent% (e.g. `bind:%prop%={...}` instead of `%prop%={...}`)
 * @param {string} parent
 * @param {string} prop
 * @param {string} child
 * @param {string} owner
 */ function ownership_invalid_binding(parent, prop, child, owner) {
    if (DEV) {
        console.warn(`%c[svelte] ownership_invalid_binding\n%c${parent} passed property \`${prop}\` to ${child} with \`bind:\`, but its parent component ${owner} did not declare \`${prop}\` as a binding. Consider creating a binding between ${owner} and ${parent} (e.g. \`bind:${prop}={...}\` instead of \`${prop}={...}\`)\nhttps://svelte.dev/e/ownership_invalid_binding`, bold, normal);
    } else {
        console.warn(`https://svelte.dev/e/ownership_invalid_binding`);
    }
}
/**
 * Mutating unbound props (`%name%`, at %location%) is strongly discouraged. Consider using `bind:%prop%={...}` in %parent% (or using a callback) instead
 * @param {string} name
 * @param {string} location
 * @param {string} prop
 * @param {string} parent
 */ function ownership_invalid_mutation(name, location, prop, parent) {
    if (DEV) {
        console.warn(`%c[svelte] ownership_invalid_mutation\n%cMutating unbound props (\`${name}\`, at ${location}) is strongly discouraged. Consider using \`bind:${prop}={...}\` in ${parent} (or using a callback) instead\nhttps://svelte.dev/e/ownership_invalid_mutation`, bold, normal);
    } else {
        console.warn(`https://svelte.dev/e/ownership_invalid_mutation`);
    }
}
/**
 * The `value` property of a `<select multiple>` element should be an array, but it received a non-array value. The selection will be kept as is.
 */ function select_multiple_invalid_value() {
    if (DEV) {
        console.warn(`%c[svelte] select_multiple_invalid_value\n%cThe \`value\` property of a \`<select multiple>\` element should be an array, but it received a non-array value. The selection will be kept as is.\nhttps://svelte.dev/e/select_multiple_invalid_value`, bold, normal);
    } else {
        console.warn(`https://svelte.dev/e/select_multiple_invalid_value`);
    }
}
/**
 * Reactive `$state(...)` proxies and the values they proxy have different identities. Because of this, comparisons with `%operator%` will produce unexpected results
 * @param {string} operator
 */ function state_proxy_equality_mismatch(operator) {
    if (esm_env__WEBPACK_IMPORTED_MODULE_0__/* ["default"] */.A) {
        console.warn(`%c[svelte] state_proxy_equality_mismatch\n%cReactive \`$state(...)\` proxies and the values they proxy have different identities. Because of this, comparisons with \`${operator}\` will produce unexpected results\nhttps://svelte.dev/e/state_proxy_equality_mismatch`, bold, normal);
    } else {
        console.warn(`https://svelte.dev/e/state_proxy_equality_mismatch`);
    }
}
/**
 * Tried to unmount a state proxy, rather than a component
 */ function state_proxy_unmount() {
    if (esm_env__WEBPACK_IMPORTED_MODULE_0__/* ["default"] */.A) {
        console.warn(`%c[svelte] state_proxy_unmount\n%cTried to unmount a state proxy, rather than a component\nhttps://svelte.dev/e/state_proxy_unmount`, bold, normal);
    } else {
        console.warn(`https://svelte.dev/e/state_proxy_unmount`);
    }
}
/**
 * A `<svelte:boundary>` `reset` function only resets the boundary the first time it is called
 */ function svelte_boundary_reset_noop() {
    if (esm_env__WEBPACK_IMPORTED_MODULE_0__/* ["default"] */.A) {
        console.warn(`%c[svelte] svelte_boundary_reset_noop\n%cA \`<svelte:boundary>\` \`reset\` function only resets the boundary the first time it is called\nhttps://svelte.dev/e/svelte_boundary_reset_noop`, bold, normal);
    } else {
        console.warn(`https://svelte.dev/e/svelte_boundary_reset_noop`);
    }
}
/**
 * The `slide` transition does not work correctly for elements with `display: %value%`
 * @param {string} value
 */ function transition_slide_display(value) {
    if (DEV) {
        console.warn(`%c[svelte] transition_slide_display\n%cThe \`slide\` transition does not work correctly for elements with \`display: ${value}\`\nhttps://svelte.dev/e/transition_slide_display`, bold, normal);
    } else {
        console.warn(`https://svelte.dev/e/transition_slide_display`);
    }
}


}),
999: (function () {

;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/version.js
// generated during release, do not modify
/**
 * The current version, as set in package.json.
 * @type {string}
 */ const VERSION = '5.39.3';
const PUBLIC_VERSION = '5';

;// CONCATENATED MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/disclose-version.js

if (typeof window !== 'undefined') {
    var _ref, _window;
    // @ts-expect-error
    ((_ref = (_window = window).__svelte ?? (_window.__svelte = {})).v ?? (_ref.v = new Set())).add(PUBLIC_VERSION);
}


}),
817: (function (__unused_webpack_module, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  I0: () => (async_mode_flag),
  LM: () => (legacy_mode_flag),
  Ny: () => (enable_legacy_mode_flag),
  _G: () => (tracing_mode_flag)
});
let async_mode_flag = false;
let legacy_mode_flag = false;
let tracing_mode_flag = false;
function enable_async_mode_flag() {
    async_mode_flag = true;
}
/** ONLY USE THIS DURING TESTING */ function disable_async_mode_flag() {
    async_mode_flag = false;
}
function enable_legacy_mode_flag() {
    legacy_mode_flag = true;
}
function enable_tracing_mode_flag() {
    tracing_mode_flag = true;
}


}),
306: (function (__unused_webpack_module, __unused_webpack___webpack_exports__, __webpack_require__) {
/* ESM import */var _index_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(817);

(0,_index_js__WEBPACK_IMPORTED_MODULE_0__/* .enable_legacy_mode_flag */.Ny)();


}),
826: (function (__unused_webpack_module, __unused_webpack___webpack_exports__, __webpack_require__) {
/* ESM import */var _utils_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(986);
/** @import { Snapshot } from './types' */ 


/**
 * In dev, we keep track of which properties could not be cloned. In prod
 * we don't bother, but we keep a dummy array around so that the
 * signature stays the same
 * @type {string[]}
 */ const empty = (/* unused pure expression or super */ null && ([]));
/**
 * @template T
 * @param {T} value
 * @param {boolean} [skip_warning]
 * @param {boolean} [no_tojson]
 * @returns {Snapshot<T>}
 */ function snapshot(value) {
    let skip_warning = arguments.length > 1 && arguments[1] !== void 0 ? arguments[1] : false, no_tojson = arguments.length > 2 && arguments[2] !== void 0 ? arguments[2] : false;
    if (DEV && !skip_warning) {
        /** @type {string[]} */ const paths = [];
        const copy = clone(value, new Map(), '', paths, null, no_tojson);
        if (paths.length === 1 && paths[0] === '') {
            // value could not be cloned
            w.state_snapshot_uncloneable();
        } else if (paths.length > 0) {
            // some properties could not be cloned
            const slice = paths.length > 10 ? paths.slice(0, 7) : paths.slice(0, 10);
            const excess = paths.length - slice.length;
            let uncloned = slice.map((path)=>`- <value>${path}`).join('\n');
            if (excess > 0) uncloned += `\n- ...and ${excess} more`;
            w.state_snapshot_uncloneable(uncloned);
        }
        return copy;
    }
    return clone(value, new Map(), '', empty, null, no_tojson);
}
/**
 * @template T
 * @param {T} value
 * @param {Map<T, Snapshot<T>>} cloned
 * @param {string} path
 * @param {string[]} paths
 * @param {null | T} [original] The original value, if `value` was produced from a `toJSON` call
 * @param {boolean} [no_tojson]
 * @returns {Snapshot<T>}
 */ function clone(value, cloned, path, paths) {
    let original = arguments.length > 4 && arguments[4] !== void 0 ? arguments[4] : null, no_tojson = arguments.length > 5 && arguments[5] !== void 0 ? arguments[5] : false;
    if (typeof value === 'object' && value !== null) {
        var unwrapped = cloned.get(value);
        if (unwrapped !== undefined) return unwrapped;
        if (value instanceof Map) return /** @type {Snapshot<T>} */ new Map(value);
        if (value instanceof Set) return /** @type {Snapshot<T>} */ new Set(value);
        if (is_array(value)) {
            var copy = /** @type {Snapshot<any>} */ Array(value.length);
            cloned.set(value, copy);
            if (original !== null) {
                cloned.set(original, copy);
            }
            for(var i = 0; i < value.length; i += 1){
                var element = value[i];
                if (i in value) {
                    copy[i] = clone(element, cloned, DEV ? `${path}[${i}]` : path, paths, null, no_tojson);
                }
            }
            return copy;
        }
        if (get_prototype_of(value) === object_prototype) {
            /** @type {Snapshot<any>} */ copy = {};
            cloned.set(value, copy);
            if (original !== null) {
                cloned.set(original, copy);
            }
            for(var key in value){
                copy[key] = clone(// @ts-expect-error
                value[key], cloned, DEV ? `${path}.${key}` : path, paths, null, no_tojson);
            }
            return copy;
        }
        if (value instanceof Date) {
            return /** @type {Snapshot<T>} */ structuredClone(value);
        }
        if (typeof /** @type {T & { toJSON?: any } } */ value.toJSON === 'function' && !no_tojson) {
            return clone(/** @type {T & { toJSON(): any } } */ value.toJSON(), cloned, DEV ? `${path}.toJSON()` : path, paths, // Associate the instance with the toJSON clone
            value);
        }
    }
    if (value instanceof EventTarget) {
        // can't be cloned
        return /** @type {Snapshot<T>} */ value;
    }
    try {
        return /** @type {Snapshot<T>} */ structuredClone(value);
    } catch (e) {
        if (DEV) {
            paths.push(path);
        }
        return /** @type {Snapshot<T>} */ value;
    }
}


}),
986: (function (__unused_webpack_module, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  CL: () => (get_descriptors),
  J8: () => (get_descriptor),
  N7: () => (object_prototype),
  Oh: () => (get_prototype_of),
  PI: () => (is_array),
  Qu: () => (define_property),
  ZZ: () => (is_extensible),
  bg: () => (array_from),
  d$: () => (object_keys),
  lQ: () => (noop),
  lc: () => (index_of),
  oO: () => (run_all),
  ve: () => (array_prototype),
  yX: () => (deferred)
});
// Store the references to globals in case someone tries to monkey patch these, causing the below
// to de-opt (this occurs often when using popular extensions).
var is_array = Array.isArray;
var index_of = Array.prototype.indexOf;
var array_from = Array.from;
var object_keys = Object.keys;
var define_property = Object.defineProperty;
var get_descriptor = Object.getOwnPropertyDescriptor;
var get_descriptors = Object.getOwnPropertyDescriptors;
var object_prototype = Object.prototype;
var array_prototype = Array.prototype;
var get_prototype_of = Object.getPrototypeOf;
var is_extensible = Object.isExtensible;
/**
 * @param {any} thing
 * @returns {thing is Function}
 */ function is_function(thing) {
    return typeof thing === 'function';
}
const noop = ()=>{};
// Adapted from https://github.com/then/is-promise/blob/master/index.js
// Distributed under MIT License https://github.com/then/is-promise/blob/master/LICENSE
/**
 * @template [T=any]
 * @param {any} value
 * @returns {value is PromiseLike<T>}
 */ function is_promise(value) {
    return typeof (value === null || value === void 0 ? void 0 : value.then) === 'function';
}
/** @param {Function} fn */ function run(fn) {
    return fn();
}
/** @param {Array<() => void>} arr */ function run_all(arr) {
    for(var i = 0; i < arr.length; i++){
        arr[i]();
    }
}
/**
 * TODO replace with Promise.withResolvers once supported widely enough
 * @template T
 */ function deferred() {
    /** @type {(value: T) => void} */ var resolve;
    /** @type {(reason: any) => void} */ var reject;
    /** @type {Promise<T>} */ var promise = new Promise((res, rej)=>{
        resolve = res;
        reject = rej;
    });
    // @ts-expect-error
    return {
        promise,
        resolve,
        reject
    };
}
/**
 * @template V
 * @param {V} value
 * @param {V | (() => V)} fallback
 * @param {boolean} [lazy]
 * @returns {V}
 */ function fallback(value, fallback) {
    let lazy = arguments.length > 2 && arguments[2] !== void 0 ? arguments[2] : false;
    return value === undefined ? lazy ? /** @type {() => V} */ fallback() : /** @type {V} */ fallback : value;
}
/**
 * When encountering a situation like `let [a, b, c] = $derived(blah())`,
 * we need to stash an intermediate value that `a`, `b`, and `c` derive
 * from, in case it's an iterable
 * @template T
 * @param {ArrayLike<T> | Iterable<T>} value
 * @param {number} [n]
 * @returns {Array<T>}
 */ function to_array(value, n) {
    // return arrays unchanged
    if (Array.isArray(value)) {
        return value;
    }
    // if value is not iterable, or `n` is unspecified (indicates a rest
    // element, which means we're not concerned about unbounded iterables)
    // convert to an array with `Array.from`
    if (n === undefined || !(Symbol.iterator in value)) {
        return Array.from(value);
    }
    // otherwise, populate an array with `n` values
    /** @type {T[]} */ const array = [];
    for (const element of value){
        array.push(element);
        if (array.length === n) break;
    }
    return array;
}


}),
461: (function (__unused_webpack_module, __unused_webpack___webpack_exports__, __webpack_require__) {
/* ESM import */var _utils_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(314);




/**
 * @param {() => string} tag_fn
 * @returns {void}
 */ function validate_void_dynamic_element(tag_fn) {
    const tag = tag_fn();
    if (tag && is_void(tag)) {
        w.dynamic_void_element_content(tag);
    }
}
/** @param {() => unknown} tag_fn */ function validate_dynamic_element_tag(tag_fn) {
    const tag = tag_fn();
    const is_string = typeof tag === 'string';
    if (tag && !is_string) {
        e.svelte_element_invalid_this_value();
    }
}
/**
 * @param {any} store
 * @param {string} name
 */ function validate_store(store, name) {
    if (store != null && typeof store.subscribe !== 'function') {
        e.store_invalid_shape(name);
    }
}
/**
 * @template {(...args: any[]) => unknown} T
 * @param {T} fn
 */ function prevent_snippet_stringification(fn) {
    fn.toString = ()=>{
        e.snippet_without_render_tag();
        return '';
    };
    return fn;
}


}),
314: (function (__unused_webpack_module, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  GY: () => (is_passive_event),
  If: () => (sanitize_location),
  tW: () => (hash)
});
const regex_return_characters = /\r/g;
/**
 * @param {string} str
 * @returns {string}
 */ function hash(str) {
    str = str.replace(regex_return_characters, '');
    let hash = 5381;
    let i = str.length;
    while(i--)hash = (hash << 5) - hash ^ str.charCodeAt(i);
    return (hash >>> 0).toString(36);
}
const VOID_ELEMENT_NAMES = (/* unused pure expression or super */ null && ([
    'area',
    'base',
    'br',
    'col',
    'command',
    'embed',
    'hr',
    'img',
    'input',
    'keygen',
    'link',
    'meta',
    'param',
    'source',
    'track',
    'wbr'
]));
/**
 * Returns `true` if `name` is of a void element
 * @param {string} name
 */ function is_void(name) {
    return VOID_ELEMENT_NAMES.includes(name) || name.toLowerCase() === '!doctype';
}
const RESERVED_WORDS = (/* unused pure expression or super */ null && ([
    'arguments',
    'await',
    'break',
    'case',
    'catch',
    'class',
    'const',
    'continue',
    'debugger',
    'default',
    'delete',
    'do',
    'else',
    'enum',
    'eval',
    'export',
    'extends',
    'false',
    'finally',
    'for',
    'function',
    'if',
    'implements',
    'import',
    'in',
    'instanceof',
    'interface',
    'let',
    'new',
    'null',
    'package',
    'private',
    'protected',
    'public',
    'return',
    'static',
    'super',
    'switch',
    'this',
    'throw',
    'true',
    'try',
    'typeof',
    'var',
    'void',
    'while',
    'with',
    'yield'
]));
/**
 * Returns `true` if `word` is a reserved JavaScript keyword
 * @param {string} word
 */ function is_reserved(word) {
    return RESERVED_WORDS.includes(word);
}
/**
 * @param {string} name
 */ function is_capture_event(name) {
    return name.endsWith('capture') && name !== 'gotpointercapture' && name !== 'lostpointercapture';
}
/** List of Element events that will be delegated */ const DELEGATED_EVENTS = (/* unused pure expression or super */ null && ([
    'beforeinput',
    'click',
    'change',
    'dblclick',
    'contextmenu',
    'focusin',
    'focusout',
    'input',
    'keydown',
    'keyup',
    'mousedown',
    'mousemove',
    'mouseout',
    'mouseover',
    'mouseup',
    'pointerdown',
    'pointermove',
    'pointerout',
    'pointerover',
    'pointerup',
    'touchend',
    'touchmove',
    'touchstart'
]));
/**
 * Returns `true` if `event_name` is a delegated event
 * @param {string} event_name
 */ function is_delegated(event_name) {
    return DELEGATED_EVENTS.includes(event_name);
}
/**
 * Attributes that are boolean, i.e. they are present or not present.
 */ const DOM_BOOLEAN_ATTRIBUTES = [
    'allowfullscreen',
    'async',
    'autofocus',
    'autoplay',
    'checked',
    'controls',
    'default',
    'disabled',
    'formnovalidate',
    'indeterminate',
    'inert',
    'ismap',
    'loop',
    'multiple',
    'muted',
    'nomodule',
    'novalidate',
    'open',
    'playsinline',
    'readonly',
    'required',
    'reversed',
    'seamless',
    'selected',
    'webkitdirectory',
    'defer',
    'disablepictureinpicture',
    'disableremoteplayback'
];
/**
 * Returns `true` if `name` is a boolean attribute
 * @param {string} name
 */ function is_boolean_attribute(name) {
    return DOM_BOOLEAN_ATTRIBUTES.includes(name);
}
/**
 * @type {Record<string, string>}
 * List of attribute names that should be aliased to their property names
 * because they behave differently between setting them as an attribute and
 * setting them as a property.
 */ const ATTRIBUTE_ALIASES = (/* unused pure expression or super */ null && ({
    // no `class: 'className'` because we handle that separately
    formnovalidate: 'formNoValidate',
    ismap: 'isMap',
    nomodule: 'noModule',
    playsinline: 'playsInline',
    readonly: 'readOnly',
    defaultvalue: 'defaultValue',
    defaultchecked: 'defaultChecked',
    srcobject: 'srcObject',
    novalidate: 'noValidate',
    allowfullscreen: 'allowFullscreen',
    disablepictureinpicture: 'disablePictureInPicture',
    disableremoteplayback: 'disableRemotePlayback'
}));
/**
 * @param {string} name
 */ function normalize_attribute(name) {
    name = name.toLowerCase();
    return ATTRIBUTE_ALIASES[name] ?? name;
}
const DOM_PROPERTIES = [
    ...DOM_BOOLEAN_ATTRIBUTES,
    'formNoValidate',
    'isMap',
    'noModule',
    'playsInline',
    'readOnly',
    'value',
    'volume',
    'defaultValue',
    'defaultChecked',
    'srcObject',
    'noValidate',
    'allowFullscreen',
    'disablePictureInPicture',
    'disableRemotePlayback'
];
/**
 * @param {string} name
 */ function is_dom_property(name) {
    return DOM_PROPERTIES.includes(name);
}
const NON_STATIC_PROPERTIES = (/* unused pure expression or super */ null && ([
    'autofocus',
    'muted',
    'defaultValue',
    'defaultChecked'
]));
/**
 * Returns `true` if the given attribute cannot be set through the template
 * string, i.e. needs some kind of JavaScript handling to work.
 * @param {string} name
 */ function cannot_be_set_statically(name) {
    return NON_STATIC_PROPERTIES.includes(name);
}
/**
 * Subset of delegated events which should be passive by default.
 * These two are already passive via browser defaults on window, document and body.
 * But since
 * - we're delegating them
 * - they happen often
 * - they apply to mobile which is generally less performant
 * we're marking them as passive by default for other elements, too.
 */ const PASSIVE_EVENTS = [
    'touchstart',
    'touchmove'
];
/**
 * Returns `true` if `name` is a passive event
 * @param {string} name
 */ function is_passive_event(name) {
    return PASSIVE_EVENTS.includes(name);
}
const CONTENT_EDITABLE_BINDINGS = (/* unused pure expression or super */ null && ([
    'textContent',
    'innerHTML',
    'innerText'
]));
/** @param {string} name */ function is_content_editable_binding(name) {
    return CONTENT_EDITABLE_BINDINGS.includes(name);
}
const LOAD_ERROR_ELEMENTS = (/* unused pure expression or super */ null && ([
    'body',
    'embed',
    'iframe',
    'img',
    'link',
    'object',
    'script',
    'style',
    'track'
]));
/**
 * Returns `true` if the element emits `load` and `error` events
 * @param {string} name
 */ function is_load_error_element(name) {
    return LOAD_ERROR_ELEMENTS.includes(name);
}
const SVG_ELEMENTS = (/* unused pure expression or super */ null && ([
    'altGlyph',
    'altGlyphDef',
    'altGlyphItem',
    'animate',
    'animateColor',
    'animateMotion',
    'animateTransform',
    'circle',
    'clipPath',
    'color-profile',
    'cursor',
    'defs',
    'desc',
    'discard',
    'ellipse',
    'feBlend',
    'feColorMatrix',
    'feComponentTransfer',
    'feComposite',
    'feConvolveMatrix',
    'feDiffuseLighting',
    'feDisplacementMap',
    'feDistantLight',
    'feDropShadow',
    'feFlood',
    'feFuncA',
    'feFuncB',
    'feFuncG',
    'feFuncR',
    'feGaussianBlur',
    'feImage',
    'feMerge',
    'feMergeNode',
    'feMorphology',
    'feOffset',
    'fePointLight',
    'feSpecularLighting',
    'feSpotLight',
    'feTile',
    'feTurbulence',
    'filter',
    'font',
    'font-face',
    'font-face-format',
    'font-face-name',
    'font-face-src',
    'font-face-uri',
    'foreignObject',
    'g',
    'glyph',
    'glyphRef',
    'hatch',
    'hatchpath',
    'hkern',
    'image',
    'line',
    'linearGradient',
    'marker',
    'mask',
    'mesh',
    'meshgradient',
    'meshpatch',
    'meshrow',
    'metadata',
    'missing-glyph',
    'mpath',
    'path',
    'pattern',
    'polygon',
    'polyline',
    'radialGradient',
    'rect',
    'set',
    'solidcolor',
    'stop',
    'svg',
    'switch',
    'symbol',
    'text',
    'textPath',
    'tref',
    'tspan',
    'unknown',
    'use',
    'view',
    'vkern'
]));
/** @param {string} name */ function is_svg(name) {
    return SVG_ELEMENTS.includes(name);
}
const MATHML_ELEMENTS = (/* unused pure expression or super */ null && ([
    'annotation',
    'annotation-xml',
    'maction',
    'math',
    'merror',
    'mfrac',
    'mi',
    'mmultiscripts',
    'mn',
    'mo',
    'mover',
    'mpadded',
    'mphantom',
    'mprescripts',
    'mroot',
    'mrow',
    'ms',
    'mspace',
    'msqrt',
    'mstyle',
    'msub',
    'msubsup',
    'msup',
    'mtable',
    'mtd',
    'mtext',
    'mtr',
    'munder',
    'munderover',
    'semantics'
]));
/** @param {string} name */ function is_mathml(name) {
    return MATHML_ELEMENTS.includes(name);
}
const STATE_CREATION_RUNES = /** @type {const} */ [
    '$state',
    '$state.raw',
    '$derived',
    '$derived.by'
];
const RUNES = /** @type {const} */ [
    ...STATE_CREATION_RUNES,
    '$state.snapshot',
    '$props',
    '$props.id',
    '$bindable',
    '$effect',
    '$effect.pre',
    '$effect.tracking',
    '$effect.root',
    '$effect.pending',
    '$inspect',
    '$inspect().with',
    '$inspect.trace',
    '$host'
];
/** @typedef {typeof RUNES[number]} RuneName */ /**
 * @param {string} name
 * @returns {name is RuneName}
 */ function is_rune(name) {
    return RUNES.includes(/** @type {RuneName} */ name);
}
/** @typedef {typeof STATE_CREATION_RUNES[number]} StateCreationRuneName */ /**
 * @param {string} name
 * @returns {name is StateCreationRuneName}
 */ function is_state_creation_rune(name) {
    return STATE_CREATION_RUNES.includes(/** @type {StateCreationRuneName} */ name);
}
/** List of elements that require raw contents and should not have SSR comments put in them */ const RAW_TEXT_ELEMENTS = /** @type {const} */ (/* unused pure expression or super */ null && ([
    'textarea',
    'script',
    'style',
    'title'
]));
/** @param {string} name */ function is_raw_text_element(name) {
    return RAW_TEXT_ELEMENTS.includes(/** @type {typeof RAW_TEXT_ELEMENTS[number]} */ name);
}
/**
 * Prevent devtools trying to make `location` a clickable link by inserting a zero-width space
 * @template {string | undefined} T
 * @param {T} location
 * @returns {T};
 */ function sanitize_location(location) {
    return /** @type {T} */ location === null || location === void 0 ? void 0 : location.replace(/\//g, '/\u200b');
}


}),
396: (function (__unused_webpack___webpack_module__, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  _: () => (_check_private_redeclaration)
});
function _check_private_redeclaration(obj, privateCollection) {
    if (privateCollection.has(obj)) {
        throw new TypeError("Cannot initialize the same private elements twice on an object");
    }
}



}),
809: (function (__unused_webpack___webpack_module__, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  _: () => (_class_extract_field_descriptor)
});
function _class_extract_field_descriptor(receiver, privateMap, action) {
    if (!privateMap.has(receiver)) throw new TypeError("attempted to " + action + " private field on non-instance");

    return privateMap.get(receiver);
}



}),
570: (function (__unused_webpack___webpack_module__, __webpack_exports__, __webpack_require__) {

// EXPORTS
__webpack_require__.d(__webpack_exports__, {
  _: () => (/* binding */ _class_private_field_get)
});

;// CONCATENATED MODULE: ./node_modules/.pnpm/@swc+helpers@0.5.17/node_modules/@swc/helpers/esm/_class_apply_descriptor_get.js
function _class_apply_descriptor_get(receiver, descriptor) {
    if (descriptor.get) return descriptor.get.call(receiver);

    return descriptor.value;
}


// EXTERNAL MODULE: ./node_modules/.pnpm/@swc+helpers@0.5.17/node_modules/@swc/helpers/esm/_class_extract_field_descriptor.js
var _class_extract_field_descriptor = __webpack_require__(809);
;// CONCATENATED MODULE: ./node_modules/.pnpm/@swc+helpers@0.5.17/node_modules/@swc/helpers/esm/_class_private_field_get.js



function _class_private_field_get(receiver, privateMap) {
    var descriptor = (0,_class_extract_field_descriptor._)(receiver, privateMap, "get");
    return _class_apply_descriptor_get(receiver, descriptor);
}



}),
636: (function (__unused_webpack___webpack_module__, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  _: () => (_class_private_field_init)
});
/* ESM import */var _check_private_redeclaration_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(396);


function _class_private_field_init(obj, privateMap, value) {
    (0,_check_private_redeclaration_js__WEBPACK_IMPORTED_MODULE_0__._)(obj, privateMap);
    privateMap.set(obj, value);
}



}),
549: (function (__unused_webpack___webpack_module__, __webpack_exports__, __webpack_require__) {

// EXPORTS
__webpack_require__.d(__webpack_exports__, {
  _: () => (/* binding */ _class_private_field_set)
});

;// CONCATENATED MODULE: ./node_modules/.pnpm/@swc+helpers@0.5.17/node_modules/@swc/helpers/esm/_class_apply_descriptor_set.js
function _class_apply_descriptor_set(receiver, descriptor, value) {
    if (descriptor.set) descriptor.set.call(receiver, value);
    else {
        if (!descriptor.writable) {
            // This should only throw in strict mode, but class bodies are
            // always strict and private fields can only be used inside
            // class bodies.
            throw new TypeError("attempted to set read only private field");
        }
        descriptor.value = value;
    }
}


// EXTERNAL MODULE: ./node_modules/.pnpm/@swc+helpers@0.5.17/node_modules/@swc/helpers/esm/_class_extract_field_descriptor.js
var _class_extract_field_descriptor = __webpack_require__(809);
;// CONCATENATED MODULE: ./node_modules/.pnpm/@swc+helpers@0.5.17/node_modules/@swc/helpers/esm/_class_private_field_set.js



function _class_private_field_set(receiver, privateMap, value) {
    var descriptor = (0,_class_extract_field_descriptor._)(receiver, privateMap, "set");
    _class_apply_descriptor_set(receiver, descriptor, value);
    return value;
}



}),
585: (function (__unused_webpack___webpack_module__, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  _: () => (_class_private_method_get)
});
function _class_private_method_get(receiver, privateSet, fn) {
    if (!privateSet.has(receiver)) throw new TypeError("attempted to get private field on non-instance");

    return fn;
}



}),
23: (function (__unused_webpack___webpack_module__, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  _: () => (_class_private_method_init)
});
/* ESM import */var _check_private_redeclaration_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(396);


function _class_private_method_init(obj, privateSet) {
    (0,_check_private_redeclaration_js__WEBPACK_IMPORTED_MODULE_0__._)(obj, privateSet);
    privateSet.add(obj);
}



}),
925: (function (__unused_webpack___webpack_module__, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  _: () => (_define_property)
});
function _define_property(obj, key, value) {
    if (key in obj) {
        Object.defineProperty(obj, key, { value: value, enumerable: true, configurable: true, writable: true });
    } else obj[key] = value;

    return obj;
}



}),
832: (function (__unused_webpack___webpack_module__, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  A: () => (__WEBPACK_DEFAULT_EXPORT__)
});
/* ESM default export */ const __WEBPACK_DEFAULT_EXPORT__ = (false);


}),
521: (function (__unused_webpack___webpack_module__, __webpack_exports__, __webpack_require__) {
__webpack_require__.d(__webpack_exports__, {
  w: () => (route)
});
function navigate(path) {
    history.pushState(null, "", path);
    window.dispatchEvent(new CustomEvent('svelteNavigate', { detail: { path } }));
}
function route(node) {
    const handleClick = (event) => {
        // Only handle if it's a left-click without modifier keys
        if (event.button === 0 && !event.ctrlKey && !event.metaKey && !event.altKey && !event.shiftKey) {
            event.preventDefault();
            // Get the href from the anchor tag
            const href = node.getAttribute('href') || '';
            // Handle external links
            const isExternal = href.startsWith('http://') ||
                href.startsWith('https://') ||
                href.startsWith('//') ||
                node.hasAttribute('external');
            if (!isExternal) {
                // Determine the path based on whether it's absolute or relative
                let path;
                if (href.startsWith('/')) {
                    // Absolute path - use as is
                    path = href;
                }
                else if (href === '' || href === '#') {
                    // Empty href or hash only - stay on current page
                    path = window.location.pathname;
                }
                else {
                    // Relative path - combine with current path
                    const currentPath = window.location.pathname;
                    // Ensure the current path ends with a slash if not the root
                    const basePath = currentPath === '/'
                        ? '/'
                        : currentPath.endsWith('/')
                            ? currentPath
                            : currentPath + '/';
                    // Combine base path with relative href
                    path = basePath + href;
                }
                // Clean up any double slashes (except after protocol)
                path = path.replace(/([^:]\/)\/+/g, '$1');
                // Handle relative paths with ../
                if (path.includes('../')) {
                    const segments = path.split('/');
                    const cleanSegments = [];
                    for (const segment of segments) {
                        if (segment === '..') {
                            // Go up one level by removing the last segment
                            if (cleanSegments.length > 1) { // Ensure we don't go above root
                                cleanSegments.pop();
                            }
                        }
                        else if (segment !== '' && segment !== '.') {
                            // Add non-empty segments that aren't current directory
                            cleanSegments.push(segment);
                        }
                    }
                    // Reconstruct the path
                    path = '/' + cleanSegments.join('/');
                }
                // Navigate to the computed path
                navigate(path);
            }
            else {
                // For external links, just follow the href
                window.location.href = href;
            }
        }
    };
    // Add event listener
    node.addEventListener('click', handleClick);
    // Return the destroy method
    return {
        destroy() {
            node.removeEventListener('click', handleClick);
        }
    };
}


}),

}]);