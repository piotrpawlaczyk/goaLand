import { createBrowserRouter } from "react-router";
import { MainPage } from "./pages/main-page";
import { CompetitionResultsPage } from "./pages/competition-results-page";
import { CompetitionLayout } from "./layouts/competition-layout";
import { CompetitionMatchesPage } from "./pages/competition-matches-page";
import { CompetitionTablePage } from "./pages/competition-table-page";
import { CompetitionGoalScorersPage } from "./pages/competition-goal-scorers-page";
import { basePath, clubaPagePath, routePath, routes } from "./config/routes";
import { ClubPage } from "./pages/club-page";
import { ErrorPage } from "./pages/error-page";

export const router = createBrowserRouter([
  {
    path: "/",
    Component: MainPage,
  },
  {
    path: routePath(routes.table),
    Component: CompetitionTablePage,
  },
  {
    path: clubaPagePath,
    Component: ClubPage,
  },
  {
    path: basePath,
    Component: CompetitionLayout,
    children: [
      { 
        path: routes.results,
        Component: CompetitionResultsPage,
      },
      {
        path: routes.matches,
        Component: CompetitionMatchesPage,
      },
      {
        path: routes.goalScorers,
        Component: CompetitionGoalScorersPage,
      },
    ],
  },
  {
    path: "*",
    Component: ErrorPage,
  },
]);
