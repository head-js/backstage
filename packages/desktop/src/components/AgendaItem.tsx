import { useEffect, useRef } from "react";
import type { AgendaItem } from "../types/agenda";
import { agendaTitle } from "../types/agenda";
import { ScopeTag } from "./ScopeTag";

interface AgendaItemRowProps {
  item: AgendaItem;
  index: number;
  onActivate: () => void;
}

export function AgendaItemRow({ item, index, onActivate }: AgendaItemRowProps) {
  const checkboxRef = useRef<HTMLInputElement>(null);
  const isHead = index === 0;
  const isIndeterminate = item.checkboxState === "indeterminate";
  const isChecked = item.checkboxState === "checked";
  const clickable = isHead && !isChecked;

  useEffect(() => {
    const el = checkboxRef.current;
    if (el) {
      el.indeterminate = isIndeterminate;
    }
  }, [isIndeterminate]);

  function handleActivate() {
    if (!clickable) {
      return;
    }
    onActivate();
  }

  function handleKeyDown(event: React.KeyboardEvent<HTMLLIElement>) {
    if (!clickable) {
      return;
    }
    if (event.key === "Enter" || event.key === " ") {
      event.preventDefault();
      onActivate();
    }
  }

  return (
    <li
      className={[
        "flex flex-col px-4 py-3",
        clickable
          ? "cursor-pointer hover:bg-base-200"
          : "cursor-not-allowed opacity-55",
      ].join(" ")}
      role="button"
      tabIndex={clickable ? 0 : -1}
      aria-disabled={!clickable}
      onClick={handleActivate}
      onKeyDown={handleKeyDown}
    >
      <div className="flex w-full items-center gap-3">
        <input
          ref={checkboxRef}
          type="checkbox"
          className="checkbox"
          checked={isChecked}
          disabled={!clickable}
          readOnly
          aria-label={agendaTitle(item.agenda)}
        />
        <span className="min-w-0 flex-1 truncate text-sm">
          {agendaTitle(item.agenda)}
        </span>
      </div>
      {isHead && item.agenda.children.length > 0 && (
        <ul
          className="pointer-events-none mt-1 ml-8 space-y-0.5"
          aria-hidden
        >
          {item.agenda.children.map((child, childIndex) => (
            <li
              key={`${child.appId}-${child.planId}-${childIndex}`}
              className="flex items-center gap-1.5 text-xs text-base-content/60"
            >
              <span aria-hidden>-</span>
              <ScopeTag left="app" value={child.appId} />
              <ScopeTag
                left="plan"
                value={child.planId}
                leftBg="#bfd5ee"
                rightBg="#cfe7ff"
              />
              <span className="min-w-0 flex-1 truncate">{child.name}</span>
            </li>
          ))}
        </ul>
      )}
    </li>
  );
}
