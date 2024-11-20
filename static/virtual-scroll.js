import { Idiomorph } from "./idiomorph.js";

export default class VirtualScroll {
  /**
   * @param {HTMLElement} container
   * @param {{numItems: number, itemHeight: number, renderItem: (index: Number) => DocumentFragment}} config
   */
  constructor(container, config) {
    this.itemHeight = config.itemHeight;
    this.numItems = config.numItems;
    this.renderItem = config.renderItem;
    this.visibleCount = Math.ceil(container.clientHeight / this.itemHeight) + 1;
    this.itemsContainer = document.createElement("div");
    this.itemsContainer.style.position = "relative";
    this.itemsContainer.style.height = `${this.numItems * this.itemHeight}px`;
    this.offset = 0;
    this.render();
    container.appendChild(this.itemsContainer);
    container.addEventListener("scroll", () => {
      this.offset = container.scrollTop;
      this.render();
    });
  }

  render() {
    const newItemsContainer = document.createElement("div");
    const startPadding = 10;
    const endPadding = 10;
    const viewportTopIndex = Math.floor(this.offset / this.itemHeight);
    const startIndex = Math.max(0, viewportTopIndex - startPadding);
    const endIndex = Math.min(
      this.numItems,
      viewportTopIndex + this.visibleCount + endPadding,
    );
    for (let i = startIndex; i < endIndex; i++) {
      const el = document.createElement("div");
      el.id = `vscroll-${i}`;
      el.style.position = "absolute";
      el.style.top = `${i * this.itemHeight}px`;
      el.style.height = `${this.itemHeight}px`;
      el.style.width = `100%`;
      el.appendChild(this.renderItem(i));
      newItemsContainer.appendChild(el);
    }
    Idiomorph.morph(this.itemsContainer, newItemsContainer.children, {
      morphStyle: "innerHTML",
    });
  }
}
