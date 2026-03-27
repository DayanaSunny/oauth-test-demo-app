function parseScopes() {
  const params = new URLSearchParams(window.location.search);
  const values = [];

  params.getAll("scope").forEach((raw) => values.push(...raw.split(/\s+/)));
  params.getAll("scopes").forEach((raw) => values.push(...raw.split(/\s+/)));

  const cleaned = values.map((v) => v.trim()).filter(Boolean);
  return Array.from(new Set(cleaned)).sort();
}

function renderScopes(scopes) {
  const list = document.getElementById("scope-list");
  const hint = document.getElementById("hint");

  if (!list || !hint) return;

  if (scopes.length === 0) {
    hint.textContent = "No scopes were provided. Try ?scope=openid%20profile%20email";
    return;
  }

  hint.textContent = "Scopes granted:";
  scopes.forEach((scope) => {
    const li = document.createElement("li");
    const code = document.createElement("code");
    code.textContent = scope;
    li.appendChild(code);
    list.appendChild(li);
  });
}

renderScopes(parseScopes());
