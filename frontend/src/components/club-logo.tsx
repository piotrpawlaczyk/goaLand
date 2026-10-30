
type ClubLogoProps = {
  url?: string;
  size?: number;
};


export function ClubLogo({ url, size }: ClubLogoProps) {
  return url && (
      <img
        className="shrink-0 object-contain"
        style={{ width: size ?? "100%", height: size ?? "100%" }}
        src={url}
        alt=""
        onError={(event) => event.currentTarget.remove()}
      />
  );
}

