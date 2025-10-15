export function navigate(path: string): void {
    history.pushState(null, "", path);
    window.dispatchEvent(new CustomEvent('svelteNavigate', { detail: { path } }));
}
export function link(node: HTMLAnchorElement): { destroy: () => void } {
  const handleClick = (event: MouseEvent) => {
      // Only handle if it's a left-click without modifier keys
      if (event.button === 0 && !event.ctrlKey && !event.metaKey && !event.altKey && !event.shiftKey) {
          event.preventDefault();

          // Get the href from the anchor tag
          const href = node.getAttribute('href') || '';

          // Handle external links
          const isExternal =
              href.startsWith('http://') ||
              href.startsWith('https://') ||
              href.startsWith('//') ||
              node.hasAttribute('external');

          if (!isExternal) {
              // Determine the path based on whether it's absolute or relative
              let path;

              if (href.startsWith('/')) {
                  // Absolute path - use as is
                  path = href;
              } else if (href === '' || href === '#') {
                  // Empty href or hash only - stay on current page
                  path = window.location.pathname;
              } else {
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
                      } else if (segment !== '' && segment !== '.') {
                          // Add non-empty segments that aren't current directory
                          cleanSegments.push(segment);
                      }
                  }

                  // Reconstruct the path
                  path = '/' + cleanSegments.join('/');
              }

              // Navigate to the computed path
              navigate(path);
          } else {
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