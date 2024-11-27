import { ReactNode, useEffect, useRef, useState } from "react";
import { Logo } from "./Logo";
import fuzzysort from "fuzzysort";
import { DownloadSimple, ListMagnifyingGlass } from "@phosphor-icons/react";
import { downloadZip } from "client-zip";
import {
	QueryClient,
	QueryClientProvider,
	useQuery,
} from "@tanstack/react-query";

const queryClient = new QueryClient();

const API_BASE = `/api/v2`;

interface Font {
	id: string;
	name: string;
	designer: string;
	subsets: string[];
	weights: string[];
	styles: string[];
}

interface Subset {
	subset: string;
	ranges: string;
}

function generateFontFaceCSS(
	font: Font,
	subset: string,
	weight: string,
	style: string,
	unicodeRange: string,
	urlPrefix: string,
): string {
	const url = `${urlPrefix}${font.id}_${subset}_${weight}_${style}.woff2`;

	return `
@font-face {
  font-family: '${font.name}';
  font-style: ${style};
  font-weight: ${weight.split("-").join(" ")};
  src: url('${url}') format('woff2');
  unicode-range: ${unicodeRange};
}
      `.trim();
}

function useFonts() {
	return useQuery<Font[]>({
		queryKey: ["fonts"],
		queryFn: async () => {
			const response = await fetch(`${API_BASE}/fonts.json`);
			if (!response.ok) {
				throw new Error(await response.text());
			}
			return await response.json();
		},
	});
}

function useSubsets() {
	return useQuery<Subset[]>({
		queryKey: ["subsets"],
		queryFn: async () => {
			const response = await fetch(`${API_BASE}/subsets.json`);
			if (!response.ok) {
				throw new Error(await response.text());
			}
			return await response.json();
		},
	});
}

function FontCSSInjector() {
	const { data: fonts } = useFonts();
	const { data: subsets } = useSubsets();
	useEffect(() => {
		if (!fonts || !subsets) return;

		let cssContent = "";

		for (const font of fonts) {
			for (const subset of font.subsets) {
				const unicodeRange = subsets.find((s) => s.subset == subset)?.ranges;
				if (!unicodeRange) {
					throw new Error(`Subset ${subset} not found`);
				}
				for (const weight of font.weights) {
					for (const style of font.styles) {
						cssContent += generateFontFaceCSS(
							font,
							subset,
							weight,
							style,
							unicodeRange,
							`${API_BASE}/download/`,
						);
					}
				}
			}
		}

		let styleElement = document.createElement("style");
		styleElement.textContent = cssContent;
		document.head.appendChild(styleElement);

		return () => {
			if (styleElement) {
				document.head.removeChild(styleElement);
			}
		};
	}, [fonts, subsets]);

	return null;
}

function VirtualScroll<T>({
	items,
	getId,
	itemHeight,
	renderItem,
}: {
	items: T[];
	itemHeight: number;
	renderItem: (item: T, index: number) => ReactNode;
	getId: (t: T) => string;
}) {
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

function Checkbox({
	label,
	checked,
	onChange,
	disabled,
}: {
	label: string;
	checked: boolean;
	onChange?: (checked: boolean) => void;
	disabled?: boolean;
}) {
	return (
		<label className="flex gap-1.5 items-center">
			<input
				type="checkbox"
				checked={checked}
				onChange={(e) => onChange?.(e.target.checked)}
				disabled={disabled}
			/>
			{label}
		</label>
	);
}

function DownloadForm({ fontId }: { fontId: string }) {
	const { data: fonts } = useFonts();
	const { data: apiSubsets } = useSubsets();

	const [allStylesChecked, setAllStylesChecked] = useState(true);
	const [allWeightsChecked, setAllWeightsChecked] = useState(true);
	const [defaultSubsetChecked, setDefaultSubsetChecked] = useState(true);

	const [selectedWeights, setSelectedWeights] = useState<string[]>([]);
	const [selectedStyles, setSelectedStyles] = useState<string[]>([]);
	const [selectedSubsets, setSelectedSubsets] = useState<string[]>([]);

	// We take the name from the fonts array here to avoid having to wait
	// for the fetch call to the API to return to render the name
	const font = fonts?.find((f) => f.id == fontId);

	async function handleDownloadClick() {
		let fontFiles: string[] = [];

		const styles = allStylesChecked ? font!.styles : selectedStyles;
		const weights = allWeightsChecked ? font!.weights : selectedWeights;
		const subsets = defaultSubsetChecked ? ["latin"] : selectedSubsets;

		let cssOutput = "";

		for (const subset of subsets) {
			const unicodeRange = apiSubsets!.find((s) => s.subset == subset)?.ranges;
			if (!unicodeRange) {
				throw new Error(`Subset ${subset} not found`);
			}
			for (const weight of weights) {
				for (const style of styles) {
					fontFiles.push(`${fontId}_${subset}_${weight}_${style}`);
					cssOutput +=
						generateFontFaceCSS(
							font!,
							subset,
							weight,
							style,
							unicodeRange,
							"",
						) + "\n";
				}
			}
		}

		const fontBlobs = await Promise.all(
			fontFiles.map((name) => fetch(`${API_BASE}/download/${name}.woff2`)),
		);

		const cssBlob = {
			name: `${fontId}.css`,
			lastModified: new Date(),
			input: cssOutput,
		};
		const files = [...fontBlobs, cssBlob];

		const blob = await downloadZip(files).blob();

		const link = document.createElement("a");
		link.href = URL.createObjectURL(blob);
		link.download = `${font!.id}-fontdelivery.zip`;
		link.click();
		link.remove();
	}

	return (
		<div className="text-sm flex flex-col gap-2">
			<p className="font-medium pr-5">Download {font?.name ?? ""}</p>
			<div>
				<div className="text-muted-foreground mb-1">Styles</div>
				{font?.styles.length == 1 ? (
					<Checkbox label={font.styles[0]} checked={true} disabled />
				) : (
					<>
						<Checkbox
							label="All styles"
							checked={allStylesChecked}
							onChange={setAllStylesChecked}
						/>
						{!allStylesChecked &&
							font?.styles.map((style) => (
								<Checkbox
									key={style}
									label={style}
									checked={selectedStyles.includes(style)}
									onChange={(checked) =>
										setSelectedStyles(
											checked
												? [...selectedStyles, style]
												: selectedStyles.filter((x) => x !== style),
										)
									}
								/>
							))}
					</>
				)}
			</div>
			<div>
				<div className="text-muted-foreground mb-1">Weights</div>
				{font?.weights.length == 1 ? (
					<Checkbox
						label={`${font.weights[0]} ${
							font.weights[0].includes("-") ? "(variable)" : "(fixed)"
						}`}
						checked={true}
						disabled
					/>
				) : (
					<>
						<Checkbox
							label="All weights"
							checked={allWeightsChecked}
							onChange={setAllWeightsChecked}
						/>
						{!allWeightsChecked &&
							font?.weights.map((weight) => (
								<Checkbox
									key={weight}
									label={weight}
									checked={selectedWeights.includes(weight)}
									onChange={(checked) =>
										setSelectedWeights(
											checked
												? [...selectedWeights, weight]
												: selectedWeights.filter((x) => x !== weight),
										)
									}
								/>
							))}
					</>
				)}
			</div>
			<div>
				<div className="text-muted-foreground mb-1">Subsets</div>

				{font?.subsets.length == 1 ? (
					<Checkbox label={font.subsets[0]} checked={true} disabled />
				) : (
					<>
						<Checkbox
							label="Default subset (latin)"
							checked={defaultSubsetChecked}
							onChange={setDefaultSubsetChecked}
						/>
						{!defaultSubsetChecked &&
							font?.subsets.map((subset) => (
								<Checkbox
									key={subset}
									label={subset}
									checked={selectedSubsets.includes(subset)}
									onChange={(checked) =>
										setSelectedSubsets(
											checked
												? [...selectedSubsets, subset]
												: selectedSubsets.filter((x) => x !== subset),
										)
									}
								/>
							))}
					</>
				)}
			</div>
			<div className="flex justify-end pt-2">
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

function FontScroller({ filter }: { filter: string }) {
	const { data: fonts } = useFonts();
	const [currentDownloadPopover, setCurrentDownloadPopover] = useState<
		string | null
	>(null);

	// Close popover with Esc-key
	useEffect(() => {
		const handleKeyDown = (event: KeyboardEvent) => {
			if (event.key === "Escape") {
				setCurrentDownloadPopover(null);
			}
		};
		document.addEventListener("keydown", handleKeyDown);
		return () => {
			document.removeEventListener("keydown", handleKeyDown);
		};
	}, [setCurrentDownloadPopover]);

	if (!fonts) return <></>;

	const sortedResult =
		filter.length > 0
			? fuzzysort.go(filter, fonts, { key: "name" }).map((x) => x.obj)
			: fonts;

	return (
		<VirtualScroll
			items={[...sortedResult]}
			getId={(font) => font.id}
			itemHeight={180}
			renderItem={(font) => (
				<div
					key={font.id}
					className="h-[180px] container mx-auto px-4 border-b"
				>
					<div className="relative flex flex-col h-full">
						<div className="absolute w-[40px] h-[178px] top-0 right-0 overflow-gradient" />
						<div className="flex flex-row justify-between mt-6">
							<span className="font-semibold">
								{font.name}{" "}
								<span className="text-muted-foreground text-sm font-normal">
									by {font.designer}
								</span>
							</span>
							<div className="relative">
								<button
									onClick={() =>
										setCurrentDownloadPopover(
											currentDownloadPopover == font.id ? null : font.id,
										)
									}
									aria-label={`Download ${font.name}`}
									className="h-12 w-12 justify-center outline-none focus:border-blue-500 border border-2 rounded text-sm flex items-center gap-1 text-md font-medium"
								>
									<DownloadSimple size={32} />
								</button>
								{currentDownloadPopover == font.id && (
									<div className="absolute mt-2 z-50 right-0 w-64 rounded-md border border-2 bg-background p-4">
										<DownloadForm fontId={font.id} />
									</div>
								)}
							</div>
						</div>
						<div
							className="text-6xl whitespace-nowrap overflow-hidden flex-grow leading-[90px]"
							style={{ fontFamily: `'${font.name}', Tofu` }}
						>
							The quick brown fox jumps over the lazy dog
						</div>
					</div>
				</div>
			)}
		/>
	);
}

function App() {
	const [filter, setFilter] = useState("");
	return (
		<QueryClientProvider client={queryClient}>
			<div className="mx-auto h-svh flex flex-col">
				<div className="container mx-auto px-4 flex justify-between items-center py-4 border-b">
					<div className="flex items-end">
						<div className="text-2xl font-semibold">
							<Logo />
						</div>
						<h1 className="font-semibold text-md hidden md:block pl-6 tracking-tight">
							webfont download service
						</h1>
					</div>
					<label className="flex items-center gap-2">
						<ListMagnifyingGlass size={32} />
						<input
							value={filter}
							onChange={(e) => setFilter(e.target.value)}
							type="text"
							aria-label="Search"
							className="border w-44 sm:w-56 outline-none focus:border-blue-500 border-2 px-2 py-1.5 rounded bg-transparent"
						/>
					</label>
				</div>
				<div className="overflow-auto flex-grow">
					<FontScroller filter={filter} />
				</div>
			</div>
			<FontCSSInjector />
		</QueryClientProvider>
	);
}

export default App;
