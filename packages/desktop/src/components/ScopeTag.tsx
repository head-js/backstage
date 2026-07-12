interface ScopeTagProps {
  left: string;
  value: string;
  leftBg?: string;
  rightBg?: string;
}

export function ScopeTag({
  left,
  value,
  leftBg = "#cbbdee",
  rightBg = "#ddcdff",
}: ScopeTagProps) {
  return (
    <span className="badge badge-sm inline-flex items-stretch gap-0 overflow-hidden p-0">
      <span
        style={{ color: "#000", backgroundColor: leftBg }}
        className="px-1.5 py-0.5 text-[10px] font-medium leading-none"
      >
        {left}
      </span>
      <span
        style={{ color: "#000", backgroundColor: rightBg }}
        className="px-1.5 py-0.5 text-[10px] font-medium leading-none"
      >
        {value}
      </span>
    </span>
  );
}
