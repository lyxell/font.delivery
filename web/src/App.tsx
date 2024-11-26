import { useEffect, useRef, useState } from "react";
import { Logo } from "./Logo";
import fuzzysort from "fuzzysort";
import { DownloadSimple, ListMagnifyingGlass, X } from "@phosphor-icons/react";
import { downloadZip } from "client-zip";
import {
	QueryClient,
	QueryClientProvider,
	useQuery,
} from "@tanstack/react-query";

import * as React from "react";
import * as Popover from "@radix-ui/react-popover";

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

interface VirtualScrollProps<T> {
	items: T[];
	itemHeight: number;
	renderItem: (item: T, index: number) => React.ReactNode;
	getId: (t: T) => string;
}

function VirtualScroll<T>({
	items,
	getId,
	itemHeight,
	renderItem,
}: VirtualScrollProps<T>) {
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
						key={getId(font)}
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

function DownloadForm({ fontId }: { fontId: string }) {
	const { data: fonts } = useFonts();
	const { data: font } = useFont(fontId);
	const [allStyles, setAllStyles] = useState(true);
	const [allWeights, setAllWeights] = useState(true);
	const [defaultSubset, setDefaultSubset] = useState(true);
	const [selectedWeights, setSelectedWeights] = useState<string[]>([]);
	const [selectedStyles, setSelectedStyles] = useState<string[]>([]);
	const [selectedSubsets, setSelectedSubsets] = useState<string[]>([]);

	// We take the name from the fonts array here to avoid having to wait
	// for the fetch call to the API to return to render the name
	const fontName = fonts?.find((f) => f.id == fontId)?.name ?? "";

	async function handleDownloadClick() {
		let fontFiles: string[] = [];

		const styles = allStyles ? font!.styles : selectedStyles;
		const weights = allWeights ? font!.weights : selectedWeights;
		const subsets = defaultSubset ? "latin" : selectedSubsets;

		for (const subset of subsets) {
			for (const weight of weights) {
				for (const style of styles) {
					fontFiles.push(`${fontId}_${subset}_${weight}_${style}`);
				}
			}
		}

		const downloads = await Promise.all(
			fontFiles.map((name) => fetch(`/api/v1/download/${name}.woff2`)),
		);
		const blob = await downloadZip(downloads).blob();

		const link = document.createElement("a");
		link.href = URL.createObjectURL(blob);
		link.download = `${font!.id}-fontdelivery.zip`;
		link.click();
		link.remove();
	}

	return (
		<div
			className="text-sm"
			style={{
				fontFamily: "Inter, Tofu",
				display: "flex",
				flexDirection: "column",
				gap: 10,
			}}
		>
			<p className="font-medium text-[15px] pr-5">Download {fontName}</p>
			<fieldset>
				<div className="text-muted-foreground mb-1">Styles</div>
				{font?.styles.length == 1 ? (
					<label className="flex gap-1.5">
						<input checked={true} disabled type="checkbox" />
						{font.styles[0]}
					</label>
				) : (
					<>
						<label className="flex gap-1.5">
							<input
								checked={allStyles}
								onChange={(e) => setAllStyles(e.target.checked)}
								type="checkbox"
							/>
							All styles
						</label>
						{!allStyles && (
							<>
								{font?.styles.map((style) => {
									return (
										<label className="flex gap-1.5">
											<input
												checked={selectedStyles.includes(style)}
												onChange={(e) => {
													setSelectedStyles(
														e.target.checked
															? [...selectedStyles, style]
															: selectedStyles.filter((x) => x != style),
													);
												}}
												type="checkbox"
											/>
											{style}
										</label>
									);
								})}
							</>
						)}
					</>
				)}
			</fieldset>
			<fieldset>
				<div className="text-muted-foreground mb-1">Weights</div>
				{font?.weights.length == 1 ? (
					<label className="flex gap-1.5">
						<input checked={true} disabled type="checkbox" />
						{font.weights[0]}{" "}
						{font.weights[0].includes("-") ? "(variable)" : "(fixed)"}
					</label>
				) : (
					<>
						<label className="flex gap-1.5">
							<input
								checked={allWeights}
								onChange={(e) => setAllWeights(e.target.checked)}
								type="checkbox"
							/>
							All weights
						</label>
						{!allWeights && (
							<>
								{font?.weights.map((weight) => {
									return (
										<label className="flex gap-1.5">
											<input
												checked={selectedWeights.includes(weight)}
												onChange={(e) => {
													setSelectedWeights(
														e.target.checked
															? [...selectedWeights, weight]
															: selectedWeights.filter((x) => x != weight),
													);
												}}
												type="checkbox"
											/>
											{weight}
										</label>
									);
								})}
							</>
						)}
					</>
				)}
			</fieldset>
			<fieldset>
				<div className="text-muted-foreground mb-1">Subsets</div>
				<label className="flex gap-1.5">
					<input
						checked={defaultSubset}
						onChange={(e) => setDefaultSubset(e.target.checked)}
						type="checkbox"
					/>
					Default subset (latin)
				</label>
				{!defaultSubset && (
					<>
						{font?.subsets.map((subset) => {
							return (
								<label className="flex gap-1.5">
									<input
										checked={selectedSubsets.includes(subset)}
										onChange={(e) => {
											setSelectedSubsets(
												e.target.checked
													? [...selectedSubsets, subset]
													: selectedSubsets.filter((x) => x != subset),
											);
										}}
										type="checkbox"
									/>
									{subset}
								</label>
							);
						})}
					</>
				)}
			</fieldset>
			<div className="flex justify-end">
				<button
					onClick={handleDownloadClick}
					className="border border-2 p-1.5 px-3 rounded font-medium outline-none focus:border-blue-500"
				>
					Download
				</button>
			</div>
		</div>
	);
}

function Main() {
	const { data: fonts } = useFonts();
	const [filter, setFilter] = useState("");

	if (!fonts) return <></>;

	const sortedResult =
		filter.length > 0
			? fuzzysort.go(filter, fonts, { key: "name" }).map((x) => x.obj)
			: fonts;

	return (
		<div
			className="container mx-auto h-screen flex flex-col px-6"
			style={{ fontFamily: "Inter, Tofu" }}
		>
			<div className="flex justify-between items-center py-4">
				<div className="text-2xl font-semibold pr-12">
					<Logo />
				</div>
				<div className="flex items-center gap-2">
					<ListMagnifyingGlass size={32} />
					<input
						value={filter}
						onChange={(e) => setFilter(e.target.value)}
						type="text"
						aria-label="Search"
						className="border outline-none focus:border-blue-500 border-2 px-2 py-1.5 rounded bg-transparent"
					/>
				</div>
			</div>
			<div className="overflow-auto flex-grow">
				<VirtualScroll
					items={[...sortedResult]}
					getId={(font) => font.id}
					itemHeight={180}
					renderItem={(font) => (
						<div
							key={font.id}
							className="py-4 h-[179px] border-b w-full flex flex-col justify-around overflow-hidden"
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
									<Popover.Root>
										<Popover.Trigger asChild>
											<button
												aria-label={`Download ${font.name}`}
												className="h-12 w-12 justify-center outline-none focus:border-blue-500 border border-2 rounded text-sm flex items-center gap-1 text-md font-medium"
											>
												<DownloadSimple size={32} />
											</button>
										</Popover.Trigger>
										<Popover.Portal>
											<Popover.Content
												align="end"
												className="w-72 rounded-md border border-2 bg-background p-4"
												sideOffset={5}
											>
												<DownloadForm fontId={font.id} />
												<Popover.Close
													className="absolute top-3 right-4 outline-none focus:text-blue-500"
													aria-label="Close"
												>
													<X size={24} />
												</Popover.Close>
											</Popover.Content>
										</Popover.Portal>
									</Popover.Root>
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
