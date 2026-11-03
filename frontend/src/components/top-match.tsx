import type { ReactNode } from "react";
import type { MatchData } from "../config/matches";
import { createClubPageLink } from "../config/routes";
import type { Club } from "../config/club";
import { Link } from "react-router";
import { ClubLogo } from "./club-logo";

export const TOP_MATCH_WIDTH = 160;

type winner = "team1" | "team2";

type TopMatchProps = {
  competitionID: number;
  homeTeam: Club;
  awayTeam: Club;
  homeTeamString: string;
  awayTeamString: string;
  winner?: winner
  separator: string;
  children: ReactNode;
};

function TopMatch({competitionID, homeTeam, awayTeam, homeTeamString, awayTeamString, 
  winner, separator, children}: TopMatchProps) {
  return (
    <div
      className="h-[56px] shrink-0 flex flex-col items-center -space-y-[2px]"
      style={{ width: TOP_MATCH_WIDTH }}
    >
      <div className="text-secondary relative flex w-full  h-[28px] items-center justify-center rounded-[5px]  bg-sections gap-[3px]">
        <ClubLogo url={homeTeam.crestUrl} size={20} />
          <div className="  flex items-center justify-center">
            <Link to={createClubPageLink(competitionID, homeTeam)} className={winner === "team1" ? "font-bold" : undefined}>
              {homeTeamString}
            </Link>
            {separator}
            <Link to={createClubPageLink(competitionID, awayTeam)}className={winner === "team2" ? "font-bold" : undefined}>
              {awayTeamString}
            </Link>
          </div>
          <ClubLogo url={awayTeam.crestUrl} size={20} />

      </div>
      <div className="flex h-[30px] w-[80px] items-center justify-center leading-[13px] bg-sections [clip-path:polygon(0_0,100%_0,87.5%_100%,12.5%_100%)]">
        <div className="text-third text-center">{children}</div>
      </div>
    </div>
  );
}



export function UpcomingMatch({competitionID, homeTeam, awayTeam, date, time,
}: MatchData) {
  return (
    <TopMatch competitionID={competitionID} homeTeam={homeTeam} homeTeamString={homeTeam.code} separator=" - " 
      awayTeam={awayTeam} awayTeamString={awayTeam.code}>
      <div>{date}</div>
      <div>{time}</div>
    </TopMatch>
  );
}

export function LiveMatch({competitionID, homeTeam, awayTeam, homeScore, awayScore,
}: MatchData) {
  return (
    <TopMatch
      competitionID={competitionID}
      homeTeam={homeTeam}
      homeTeamString={`${homeTeam.code} ${homeScore}`}
      separator=" : "
      awayTeam={awayTeam}
      awayTeamString={`${awayScore} ${awayTeam.code}`}
    >
      <span className="text-third text-red-600">trwa</span>
    </TopMatch>
  );
}

export function FinishedMatch({competitionID, homeTeam, awayTeam, homeScore, awayScore,
}: MatchData) {
  let matchWinner: winner | undefined
  if (homeScore! > awayScore!){
    matchWinner = "team1"
  }
  if (homeScore! < awayScore!){
    matchWinner = "team2"
  }

  return (
    <TopMatch
      competitionID={competitionID}
      homeTeam={homeTeam}
      homeTeamString={`${homeTeam.code} ${homeScore}`}
      separator=" : "
      awayTeam={awayTeam}
      awayTeamString={`${awayScore} ${awayTeam.code}`}
      winner={matchWinner}
    >
      <span className="text-third"> koniec </span>
    </TopMatch>
  );
}
