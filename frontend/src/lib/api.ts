import { createConnectTransport } from "@bufbuild/connect-web";
import { PromiseClient, createPromiseClient } from "@bufbuild/connect";
import type { ServiceType } from "@bufbuild/protobuf";

const transport = createConnectTransport({
  baseUrl: "/api/",
});

// NOTE: https://connect.build/docs/web/using-clients#managing-clients-and-transports
const memo = new Map();

export function getClient<T extends ServiceType>(service: T): PromiseClient<T> {
  if (memo.has(service)) {
    console.log('AHAHAHAAHA')
    return memo.get(service)!;
  }

  const client = createPromiseClient(service, transport);
  memo.set(service, client);
  return client;
}
