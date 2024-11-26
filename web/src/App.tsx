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

function VariantSelector({ id }: { id: string }) {
	const { data: font } = useFont(id);

	let fontFiles: string[] = [];

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

	if (!font) return <></>;

	return (
		<div
			style={{
				fontFamily: "Inter, Tofu",
				display: "flex",
				flexDirection: "column",
				gap: 10,
			}}
		>
			<p className="font-medium">Download {font.name}</p>
			<fieldset>
				<label className="flex gap-1.5">
					<input checked type="checkbox" />
					All styles
				</label>
			</fieldset>
			<fieldset>
				<label className="flex gap-1.5">
					<input checked type="checkbox" />
					All weights
				</label>
			</fieldset>
			<fieldset>
				<label className="flex gap-1.5">
					<input checked type="checkbox" />
					Latin charset
				</label>
			</fieldset>
			<div className="flex justify-end">
				<button
					onClick={handleDownloadClick}
					className="border border-2 p-1.5 px-3 rounded"
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
						className="border outline-none focus:border-blue-500 border-2 border-zinc-300 px-2 py-1.5 rounded"
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
									<Popover.Root>
										<Popover.Trigger asChild>
											<button
												aria-label={`Download ${font.name}`}
												className="h-12 w-12 justify-center outline-none focus:border-blue-500 border border-2 border-zinc-300 rounded text-sm flex items-center gap-1 text-md font-medium"
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
												<VariantSelector id={font.id} />
												<Popover.Close
													className="absolute top-3 right-4 outline-none focus:text-blue-500"
													aria-label="Close"
												>
													<X size={24} />
												</Popover.Close>
												<Popover.Arrow className="fill-zinc-200" />
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
