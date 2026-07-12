import { useAgenda } from "../lib/agenda/context";
import { AgendaItemRow } from "../components/AgendaItem";

export function Agenda() {
  const { items, activateHead } = useAgenda();

  if (items.length === 0) {
    return (
      <div className="rounded-md border border-base-300 bg-base-100 px-4 py-6 text-center text-sm text-base-content/55">
        No agendas available.
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <h1 className="text-base font-semibold">Agenda</h1>
      <section className="rounded-md border border-base-300 bg-base-100">
        <ul className="divide-y divide-base-300">
          {items.map((item, index) => (
            <AgendaItemRow
              key={item.agenda.id}
              item={item}
              index={index}
              onActivate={activateHead}
            />
          ))}
        </ul>
      </section>
    </div>
  );
}
