export interface AgendaChild {
  name: string;
  appId: string;
  planId: string;
}

export interface Agenda {
  id: string;
  name: string;
  children: AgendaChild[];
}

export type CheckboxState = "unchecked" | "indeterminate" | "checked";

export interface AgendaItem {
  agenda: Agenda;
  checkboxState: CheckboxState;
}

export interface AgendaState {
  items: AgendaItem[];
}

export type AgendaAction =
  | { type: "ACTIVATE_HEAD" }
  | { type: "ROTATE_COMPLETED_HEAD" };

export function agendaTitle(agenda: Agenda): string {
  return `${agenda.id}: ${agenda.name}`;
}
