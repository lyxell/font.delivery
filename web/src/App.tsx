import { create } from "zustand";
import { useEffect, useRef, useState } from "react";
import { Logo } from "./Logo";
import { BoxArrowDown, Check } from "@phosphor-icons/react";
import clsx from "clsx";
import { downloadZip } from "client-zip";
import {
	QueryClient,
	QueryClientProvider,
	useQuery,
} from "@tanstack/react-query";
import { SortedSet } from "./SortedSet";

interface Font {
	id: string;
	name: string;
	designer: string;
	subsets: string[];
	weights: string[];
	styles: string[];
}

function useFont(id: string) {
	return useQuery<Font>({
		queryKey: ["fonts", id],
		queryFn: async () => {
			const response = await fetch(`/api/v1/fonts/${id}.json`);
			if (!response.ok) {
				throw new Error(await response.text());
			}
			return await response.json();
		},
	});
}

function useFonts() {
	return useQuery<Font[]>({
		queryKey: ["fonts"],
		queryFn: async () => {
			const response = await fetch(`/api/v1/fonts.json`);
			if (!response.ok) {
				throw new Error(await response.text());
			}
			return await response.json();
		},
	});
}

interface FontStore {
	selectedFonts: SortedSet<string>;
	selectFont: (id: string) => void;

	selectedSubsets: Record<string, SortedSet<string>>;
	toggleSubset: (fontId: string, subset: string) => void;

	selectedWeights: Record<string, SortedSet<string>>;
	toggleWeight: (fontId: string, weight: string) => void;

	selectedStyles: Record<string, SortedSet<string>>;
	toggleStyle: (fontId: string, style: string) => void;
}

const useFontStore = create<FontStore>((set) => ({
	selectedFonts: new SortedSet<string>(),
	selectFont: (id: string) =>
		set((state) => {
			const newSelected = new SortedSet(state.selectedFonts.toArray());
			if (newSelected.has(id)) {
				newSelected.delete(id);
			} else {
				newSelected.add(id);
			}
			return { selectedFonts: newSelected };
		}),

	selectedSubsets: {},
	toggleSubset: (fontId: string, subset: string) =>
		set((state) => {
			const currentSubsets = state.selectedSubsets[fontId]?.toArray() ?? [];
			const newSubsets = SortedSet.fromArray(currentSubsets);
			if (newSubsets.has(subset)) {
				newSubsets.delete(subset);
			} else {
				newSubsets.add(subset);
			}
			return {
				selectedSubsets: {
					...state.selectedSubsets,
					[fontId]: newSubsets,
				},
			};
		}),

	selectedWeights: {},
	toggleWeight: (fontId: string, weight: string) =>
		set((state) => {
			const currentWeights = state.selectedWeights[fontId]?.toArray() ?? [];
			const newWeights = SortedSet.fromArray(currentWeights);
			if (newWeights.has(weight)) {
				newWeights.delete(weight);
			} else {
				newWeights.add(weight);
			}
			return {
				selectedWeights: {
					...state.selectedWeights,
					[fontId]: newWeights,
				},
			};
		}),

	selectedStyles: {},
	toggleStyle: (fontId: string, style: string) =>
		set((state) => {
			const currentStyles = state.selectedStyles[fontId]?.toArray() ?? [];
			const newStyles = SortedSet.fromArray(currentStyles);
			if (newStyles.has(style)) {
				newStyles.delete(style);
			} else {
				newStyles.add(style);
			}
			return {
				selectedStyles: {
					...state.selectedStyles,
					[fontId]: newStyles,
				},
			};
		}),
}));

interface VirtualScrollProps {
	items: Font[];
	itemHeight: number;
	renderItem: (item: Font, index: number) => React.ReactNode;
}

function VirtualScroll({ items, itemHeight, renderItem }: VirtualScrollProps) {
	const containerRef = useRef<HTMLDivElement>(null);
	const [offset, setOffset] = useState(0);

	const visibleCount = Math.ceil(window.innerHeight / itemHeight) + 1;

	const handleScroll = () => {
		if (containerRef.current) {
			setOffset(containerRef.current.scrollTop);
		}
	};

	useEffect(() => {
		const container = containerRef.current;
		if (container) {
			container.addEventListener("scroll", handleScroll);
		}
		return () => {
			if (container) {
				container.removeEventListener("scroll", handleScroll);
			}
		};
	}, []);

	const startIndex = Math.max(0, Math.floor(offset / itemHeight) - 10);
	const endIndex = Math.min(items.length, startIndex + visibleCount + 20);

	return (
		<div ref={containerRef} style={{ overflowY: "auto", height: "100%" }}>
			<div
				style={{
					position: "relative",
					height: `${items.length * itemHeight}px`,
				}}
			>
				{items.slice(startIndex, endIndex).map((font, i) => (
					<div
						key={font.id}
						style={{
							position: "absolute",
							top: `${(i + startIndex) * itemHeight}px`,
							height: `${itemHeight}px`,
							width: "100%",
						}}
					>
						{renderItem(font, i + startIndex)}
					</div>
				))}
			</div>
		</div>
	);
}

function VariantSelector({ id }: { id: string }) {
	const { data: font } = useFont(id);
	const {
		selectedWeights,
		selectedSubsets,
		selectedStyles,
		toggleWeight,
		toggleSubset,
		toggleStyle,
	} = useFontStore();

	if (!font) return <></>;

	return (
		<div>
			<div className="font-medium mt-4">Select variants for {font.name}</div>
			<div>Weights</div>
			{font.weights.map((w) => (
				<label key={w} className="block">
					<input
						type="checkbox"
						checked={selectedWeights[id]?.has(w) || false}
						onChange={() => toggleWeight(id, w)}
					/>
					{w}
				</label>
			))}
			<div>Subsets</div>
			{font.subsets.map((s) => (
				<label key={s} className="block">
					<input
						type="checkbox"
						checked={selectedSubsets[id]?.has(s) || false}
						onChange={() => toggleSubset(id, s)}
					/>
					{s}
				</label>
			))}
			<div>Styles</div>
			{font.styles.map((s) => (
				<label key={s} className="block">
					<input
						type="checkbox"
						checked={selectedStyles[id]?.has(s) || false}
						onChange={() => toggleStyle(id, s)}
					/>
					{s}
				</label>
			))}
		</div>
	);
}

function VariantSelectors() {
	const { selectedFonts, selectedWeights, selectedSubsets, selectedStyles } =
		useFontStore();

	let fontFiles: string[] = [];

	for (const id of selectedFonts.toArray()) {
		for (const subset of selectedSubsets[id]?.toArray() ?? []) {
			for (const weight of selectedWeights[id]?.toArray() ?? []) {
				for (const style of selectedStyles[id]?.toArray() ?? []) {
					fontFiles.push(`${id}_${subset}_${weight}_${style}`);
				}
			}
		}
	}

	async function handleDownloadClick() {
		const downloads = await Promise.all(
			fontFiles.map((id) => fetch(`/api/v1/download/${id}.woff2`)),
		);
		const blob = await downloadZip(downloads).blob();

		const link = document.createElement("a");
		link.href = URL.createObjectURL(blob);
		link.download = "test.zip";
		link.click();
		link.remove();
	}

	return (
		<>
			{selectedFonts.toArray().map((id) => (
				<div key={id}>
					<VariantSelector id={id} />
				</div>
			))}
			<pre>{JSON.stringify(fontFiles, null, 2)}</pre>
			<button onClick={handleDownloadClick}>Download</button>
		</>
	);
}

function Main() {
	const { selectedFonts, selectFont } = useFontStore();
	const [downloading, setDownloading] = useState(false);
	const { data: fonts } = useFonts();

	return (
		<div
			className="container mx-auto h-screen flex flex-col px-6"
			style={{ fontFamily: "Inter, Tofu" }}
		>
			<div className="flex justify-between items-center py-4">
				<div className="text-2xl font-semibold pr-12">
					<Logo />
				</div>
				<div className="text-muted-foreground">
					<button
						onClick={() => setDownloading(true)}
						className={clsx("text-sm flex flex-col items-center gap-1", {
							"text-black": selectedFonts.size() > 0,
						})}
					>
						<BoxArrowDown size={32} />
						Download
					</button>
				</div>
			</div>
			{downloading && (
				<div>
					<h2>Select variants</h2>
					<VariantSelectors />
				</div>
			)}
			<div className="overflow-auto flex-grow">
				<VirtualScroll
					items={fonts ?? []}
					itemHeight={180}
					renderItem={(font: Font) => (
						<div
							key={font.id}
							className="py-4 h-[179px] border-b border-zinc-150 w-full flex flex-col justify-around overflow-hidden"
						>
							<div className="flex justify-between">
								<span style={{ fontSize: "16px", fontWeight: 600 }}>
									{font.name}{" "}
									<span
										className="text-muted-foreground text-sm font-normal italic"
										style={{ fontFamily: `'Vollkorn', Tofu` }}
									>
										by {font.designer}
									</span>
								</span>
								<div>
									<button
										onClick={() => selectFont(font.id)}
										className={clsx(
											"border p-1.5 px-2 rounded text-sm flex items-center gap-1",
											{
												"bg-green-600 border-green-700 text-white":
													selectedFonts.has(font.id),
											},
										)}
									>
										{selectedFonts.has(font.id) ? (
											<>
												<Check size={16} />
												Selected
											</>
										) : (
											<>Select</>
										)}
									</button>
								</div>
							</div>
							<div
								className={`text-6xl whitespace-nowrap`}
								style={{ fontFamily: `'${font.name}', Tofu` }}
							>
								The quick brown fox jumps over the lazy dog
							</div>
						</div>
					)}
				/>
			</div>
		</div>
	);
}

function App() {
	return (
		<QueryClientProvider client={new QueryClient()}>
			<Main />
		</QueryClientProvider>
	);
}

export default App;
