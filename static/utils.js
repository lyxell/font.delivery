/**
 * Create a HTML element from a HTML template.
 *
 * @param {HTMLTemplateElement | string} tmpl
 * @param {Record<string, (e: Element) => void>} record
 * @returns {DocumentFragment}
 */
export function createElement(tmpl, record) {
  const template = (() => {
    if (typeof tmpl === "string") {
      const res = document.querySelector(tmpl);
      if (!(res instanceof HTMLTemplateElement)) {
        throw new Error("Query selector did not return a HTMLTemplateElement");
      }
      return res;
    }
    return tmpl;
  })();
  const element = /** @type {typeof template.content} */ (
    template.content.cloneNode(true)
  );
  for (const k of Object.keys(record)) {
    const fn = record[k];
    const e = element.querySelector(k);
    if (e !== null) {
      fn(e);
    }
  }
  return element;
}
