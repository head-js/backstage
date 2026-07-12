import { Link } from "react-router-dom";

export function NotFound() {
  return (
    <div className="flex min-h-[360px] items-center justify-center">
      <div className="text-center">
        <div className="mb-2 text-4xl font-semibold">404</div>
        <p className="mb-4 text-sm text-base-content/60">This dashboard route does not exist.</p>
        <Link to="/dashboard" className="btn btn-primary btn-sm">
          Back to Dashboard
        </Link>
      </div>
    </div>
  );
}
