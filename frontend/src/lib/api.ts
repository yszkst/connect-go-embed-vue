import { createConnectTransport } from "@connectrpc/connect-web";
import { PromiseClient, createPromiseClient } from "@connectrpc/connect";
import type { ServiceType } from "@bufbuild/protobuf";

const transport = createConnectTransport({
  baseUrl: "/api/",
});

// NOTE: https://connectrpc.com/docs/web/using-clients/#managing-clients-and-transports
const memo = new Map();

export function getClient<T extends ServiceType>(service: T): PromiseClient<T> {
  if (memo.has(service)) {
    return memo.get(service)!;
  }

  const client = createPromiseClient(service, transport);
  memo.set(service, client);
  return client;
}
