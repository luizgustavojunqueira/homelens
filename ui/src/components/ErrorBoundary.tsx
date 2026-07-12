import { Component, type ErrorInfo, type ReactNode } from "react";

interface Props {
  children?: ReactNode;
  fallback?: ReactNode;
}

interface State {
  hasError: boolean;
  error?: Error;
}

export class ErrorBoundary extends Component<Props, State> {
  public state: State = {
    hasError: false
  };

  public static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error("Uncaught error:", error, errorInfo);
  }

  public render() {
    if (this.state.hasError) {
      if (this.props.fallback) {
        return this.props.fallback;
      }
      return (
        <div className="flex flex-col items-center justify-center h-full w-full p-8 text-center bg-(--bg)">
          <h2 className="text-2xl font-bold text-red-500 mb-4">Something went wrong.</h2>
          <p className="text-(--text-dim) mb-4">
            {this.state.error?.message || "An unexpected error occurred."}
          </p>
          <button
            onClick={() => window.location.reload()}
            className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md transition-colors cursor-pointer"
          >
            Reload Page
          </button>
        </div>
      );
    }

    return this.props.children;
  }
}
