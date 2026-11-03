import { ChevronLeft, ChevronRight } from "lucide-react";
import { useEffect, useRef } from "react";
import type { MatchData } from "../config/matches";
import { FinishedMatch, LiveMatch, TOP_MATCH_WIDTH, UpcomingMatch } from "./top-match";

const NumberMatchesToShowe = 7
const TOP_MATCH_GAP = 10;
const SCROLL_STEP = TOP_MATCH_WIDTH + TOP_MATCH_GAP;
const MATCHES_WIDTH = NumberMatchesToShowe * TOP_MATCH_WIDTH + (NumberMatchesToShowe - 1) * TOP_MATCH_GAP;
const minScrollPosition = 0;

type TopMatchesProps = {
  matches: MatchData[];
};

function TopMatchItem({ match }: { match: MatchData }) {
  if (match.status === "scheduled") {
    return ( <UpcomingMatch {...match}/> );
  }

  if (match.status === "live") {
    return ( <LiveMatch {...match}/> );
  }
  return ( <FinishedMatch {...match}/> );
}

function GetScrollPosition(scrollPosition: number, maxScrollPosition: number) {
  return Math.min(Math.max(minScrollPosition, scrollPosition), maxScrollPosition);
}


export function TopMatches({ matches }: TopMatchesProps) {
  const matchesRef = useRef<HTMLDivElement>(null);
  const scrollPosition = useRef(0);
  const maxScrollPosition = Math.max(0, (matches.length - NumberMatchesToShowe) * SCROLL_STEP );


  useEffect(() => {
    const currentIndex = matches.findIndex(
      (match) => match.status === "live" || match.status === "scheduled",
    ) - Math.floor(NumberMatchesToShowe / 2);
    scrollPosition.current = GetScrollPosition(currentIndex * SCROLL_STEP, maxScrollPosition);
    matchesRef.current?.scrollTo({ left: scrollPosition.current });
  }, [matches, maxScrollPosition]);

  function scroll(direction: -1 | 1) {
    scrollPosition.current =  GetScrollPosition(scrollPosition.current + direction * SCROLL_STEP, maxScrollPosition);
    matchesRef.current?.scrollTo({
      left: scrollPosition.current,
      behavior: "smooth",
    });
  }

  return (
    <div className="mx-auto flex w-fit max-w-full items-center gap-[10px] px-[10px] py-[15px]">
      <button
        type="button"
        aria-label="Poprzednie mecze"
        className="shrink-0 cursor-pointer text-text hover:text-green-500"
        onClick={() => scroll(-1)}
      >
        <ChevronLeft size={20} strokeWidth={2.5} />
      </button>

      <div
        ref={matchesRef}
        className="flex min-w-0 items-center overflow-hidden"
        style={{ width: MATCHES_WIDTH, gap: TOP_MATCH_GAP }}
      >
        {matches.map((match) => (
          <TopMatchItem key={match.id} match={match} />
        ))}
      </div>

      <button
        type="button"
        aria-label="Następne mecze"
        className="shrink-0 cursor-pointer text-text hover:text-green-500"
        onClick={() => scroll(1)}
      >
        <ChevronRight size={20} strokeWidth={2.5} />
      </button>
    </div>
  );
}
