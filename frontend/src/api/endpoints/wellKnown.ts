import { requestJSON } from "../client";
import type { WellKnownInfo } from "../types";

export type GetInfoOptions = {
  signal?: AbortSignal;
};

export async function getInfo({
  signal,
}: GetInfoOptions = {}): Promise<WellKnownInfo> {
  const res = await requestJSON<{ info: WellKnownInfo }>("/api/v1/info", {
    signal,
  });
  return res.info;
}
