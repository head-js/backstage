import {
  createContext,
  useContext,
  useEffect,
  useReducer,
  type ReactNode,
} from "react";
import type { AgendaItem } from "../../types/agenda";
import {
  agendaReducer,
  createInitialAgendaState,
} from "./reducer";

export interface AgendaContextValue {
  items: AgendaItem[];
  activateHead: () => void;
}

const AgendaContext = createContext<AgendaContextValue | null>(null);

export function AgendaProvider({ children }: { children: ReactNode }) {
  const [state, dispatch] = useReducer(
    agendaReducer,
    undefined,
    createInitialAgendaState,
  );

  const head = state.items[0];

  useEffect(() => {
    if (!head || head.checkboxState !== "checked") {
      return;
    }

    let outerFrame = 0;
    let innerFrame = 0;

    outerFrame = requestAnimationFrame(() => {
      innerFrame = requestAnimationFrame(() => {
        dispatch({ type: "ROTATE_COMPLETED_HEAD" });
      });
    });

    return () => {
      if (outerFrame) {
        cancelAnimationFrame(outerFrame);
      }
      if (innerFrame) {
        cancelAnimationFrame(innerFrame);
      }
    };
  }, [head]);

  const value: AgendaContextValue = {
    items: state.items,
    activateHead: () => dispatch({ type: "ACTIVATE_HEAD" }),
  };

  return (
    <AgendaContext.Provider value={value}>{children}</AgendaContext.Provider>
  );
}

export function useAgenda(): AgendaContextValue {
  const ctx = useContext(AgendaContext);
  if (!ctx) {
    throw new Error("useAgenda must be used within an AgendaProvider");
  }
  return ctx;
}
