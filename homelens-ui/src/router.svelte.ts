function createRouter() {
  let path = $state(window.location.pathname);

  window.addEventListener("popstate", () => {
    path = window.location.pathname;
  });

  function go(to: string) {
    if (to === path) return;
    history.pushState({}, "", to);
    path = to;
  }

  function match(pattern: string): Record<string, string> | null {
    const patternParts = pattern.split("/").filter(Boolean);
    const pathParts = path.split("/").filter(Boolean);
    if (patternParts.length !== pathParts.length) return null;

    const params: Record<string, string> = {};
    for (let i = 0; i < patternParts.length; i++) {
      const p = patternParts[i];
      if (p.startsWith(":")) {
        params[p.slice(1)] = decodeURIComponent(pathParts[i]);
      } else if (p !== pathParts[i]) {
        return null;
      }
    }
    return params;
  }

  return {
    get path() {
      return path;
    },
    go,
    match,
  };
}

export const router = createRouter();
