export class SortedSet<T> {
	private items: T[];
	private comparator: (a: T, b: T) => number;

	constructor(
		initialItems: T[] = [],
		comparator: (a: T, b: T) => number = SortedSet.defaultComparator,
	) {
		this.comparator = comparator;
		this.items = Array.from(new Set(initialItems)).sort(this.comparator);
	}

	static defaultComparator<T>(a: T, b: T): number {
		if (a < b) return -1;
		if (a > b) return 1;
		return 0;
	}

	add(item: T): SortedSet<T> {
		if (!this.has(item)) {
			this.items = [...this.items, item].sort(this.comparator);
		}
		return this;
	}

	delete(item: T): SortedSet<T> {
		this.items = this.items.filter((i) => i !== item);
		return this;
	}

	has(item: T): boolean {
		return this.items.includes(item);
	}

	toArray(): T[] {
		return [...this.items];
	}

	size(): number {
		return this.items.length;
	}

	clear(): SortedSet<T> {
		this.items = [];
		return this;
	}

	static fromArray<T>(
		array: T[],
		comparator?: (a: T, b: T) => number,
	): SortedSet<T> {
		return new SortedSet(array, comparator);
	}
}
