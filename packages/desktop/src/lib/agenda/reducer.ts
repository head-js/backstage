import type {
  Agenda,
  AgendaAction,
  AgendaItem,
  AgendaState,
  CheckboxState,
} from "../../types/agenda";
import agendasData from "../../agendas.json";

export function createInitialAgendaState(): AgendaState {
  const items: AgendaItem[] = (agendasData as Agenda[]).map((agenda) => ({
    agenda,
    checkboxState: "unchecked",
  }));
  return { items };
}

export function agendaReducer(
  state: AgendaState,
  action: AgendaAction,
): AgendaState {
  switch (action.type) {
    case "ACTIVATE_HEAD": {
      const [head] = state.items;
      if (!head) {
        return state;
      }
      if (head.checkboxState === "checked") {
        return state;
      }
      const nextCheckboxState: CheckboxState =
        head.checkboxState === "unchecked" ? "indeterminate" : "checked";
      const nextItems = state.items.map((item, index) =>
        index === 0 ? { ...item, checkboxState: nextCheckboxState } : item,
      );
      return { items: nextItems };
    }
    case "ROTATE_COMPLETED_HEAD": {
      const [head, ...rest] = state.items;
      if (!head || head.checkboxState !== "checked") {
        return state;
      }
      const rotated: AgendaItem = {
        ...head,
        checkboxState: "unchecked",
      };
      return { items: [...rest, rotated] };
    }
    default:
      return state;
  }
}
