import { useEffect, useRef, useState } from "react";
import logo from "./logo.svg";

interface Font {
	name: string;
	designer: string;
	id: string;
}

interface VirtualScrollProps {
	items: Font[];
	itemHeight: number;
	renderItem: (item: Font, index: number) => React.ReactNode;
}

const VirtualScroll = ({
	items,
	itemHeight,
	renderItem,
}: VirtualScrollProps) => {
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
};

const App = () => {
	const [fonts, setFonts] = useState<Font[]>([]);
	const [previewText, setPreviewText] = useState("");

	useEffect(() => {
		fetch("/api/v1/fonts.json")
			.then((res) => res.json())
			.then((data) => setFonts(data));
	}, []);

	const handlePreviewTextChange = (e: React.ChangeEvent<HTMLInputElement>) => {
		setPreviewText(e.target.value);
	};

	return (
		<div className="container mx-auto h-screen flex flex-col px-6">
			<div className="flex justify-between items-center py-4">
				<div className="text-2xl font-semibold pr-12">
					<img src={logo} className="h-10" />
				</div>
				<nav className="flex gap-6">
					<input
						id="previewTextInput"
						className="border p-2 rounded-sm w-96 bg-transparent"
						placeholder="Custom text"
						value={previewText}
						onChange={handlePreviewTextChange}
					/>
				</nav>
				<div>{fonts.length} fonts available</div>
			</div>
			<div className="overflow-auto flex-grow">
				<VirtualScroll
					items={fonts}
					itemHeight={180}
					renderItem={(font) => (
						<div
							key={font.id}
							className="py-4 h-[179px] border-b border-zinc-150 w-full flex flex-col justify-around overflow-hidden"
						>
							<div className="flex justify-between">
								<span>{`${font.name} by ${font.designer}`}</span>
							</div>
							<div className={`text-7xl whitespace-nowrap font-${font.id}`}>
								{previewText.trim() ||
									"The quick brown fox jumps over the lazy dog"}
							</div>
						</div>
					)}
				/>
			</div>
		</div>
	);
};

export default App;
