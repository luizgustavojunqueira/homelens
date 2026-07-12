import { BrowserRouter, Route, Routes } from "react-router-dom";
import { Header } from "./components/header";
import { Dashboard } from "./pages/dashboard/dashboard";
import AgentDetails from "./pages/agentDetails/agentDetails";
import { useEffect } from "react";
import { connectWS, disconnectWS } from "./services/websocket/agentSocket";
import { ToastContainer } from "react-toastify";
import Settings from "./pages/settings/Settings";
import { ErrorBoundary } from "./components/ErrorBoundary";

function NotFound() {
  return (
    <div className="flex flex-col items-center justify-center h-full text-center">
      <h2 className="text-2xl font-bold text-(--text) mb-2">404 - Not Found</h2>
      <p className="text-(--text-dim)">The page you are looking for does not exist.</p>
    </div>
  );
}

function App() {
  useEffect(() => {
    connectWS();
    return () => disconnectWS();
  }, []);

  return (
    <>
      <BrowserRouter>
        <main className="relative w-full h-screen flex flex-col overflow-hidden">
          <Header />
          <ErrorBoundary>
            <Routes>
              <Route path="/" element={<Dashboard />} />
              <Route path="/agents/:guid" element={<AgentDetails />} />
              <Route path="/settings/" element={<Settings />} />
              <Route path="*" element={<NotFound />} />
            </Routes>
          </ErrorBoundary>
        </main>
      </BrowserRouter>
      <ToastContainer
        position="top-right"
        autoClose={3000}
        hideProgressBar={true}
        newestOnTop={false}
        closeOnClick
        rtl={false}
        pauseOnFocusLoss
        draggable
        pauseOnHover
        theme="dark"
      />
    </>
  );
}

export default App;
