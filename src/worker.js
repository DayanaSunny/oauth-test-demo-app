export default {
  async fetch(request, env) {
    const url = new URL(request.url);

    // Let static assets under /assets/* bypass Worker logic when configured.
    if (url.pathname.startsWith("/assets/")) {
      return env.ASSETS.fetch(request);
    }

    if (url.pathname === "/") {
      return Response.redirect(
        `${url.origin}/oauth/callback?scope=openid%20profile`,
        302
      );
    }

    // Keep the existing OAuth callback endpoint.
    if (url.pathname === "/oauth/callback") {
      return env.ASSETS.fetch(new Request(new URL("/index.html", request.url)));
    }

    // SPA fallback (works with not_found_handling as well).
    return env.ASSETS.fetch(new Request(new URL("/index.html", request.url)));
  },
};
