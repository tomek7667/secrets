import { useState, useEffect, useCallback } from "react";

export type Route = "secrets" | "users" | "tokens" | "permissions";

const validRoutes: Route[] = ["secrets", "users", "tokens", "permissions"];

function getRouteFromHash(): Route {
	const hash = window.location.hash.slice(1) as Route;
	return validRoutes.includes(hash) ? hash : "secrets";
}

export function useRouter() {
	const [route, setRouteState] = useState<Route>(getRouteFromHash);

	useEffect(() => {
		const handleHashChange = () => {
			setRouteState(getRouteFromHash());
		};

		window.addEventListener("hashchange", handleHashChange);
		return () => window.removeEventListener("hashchange", handleHashChange);
	}, []);

	const setRoute = useCallback((newRoute: Route) => {
		window.location.hash = newRoute;
	}, []);

	return { route, setRoute };
}
