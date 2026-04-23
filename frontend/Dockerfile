FROM node:22-alpine AS builder

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
ENV NUXT_TELEMETRY_DISABLED=1

WORKDIR /app

RUN corepack enable

COPY --link package.json pnpm-lock.yaml ./
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile

COPY --link . .
RUN --mount=type=cache,id=nuxt,target=/app/node_modules/.cache \
    pnpm exec nuxi build

FROM node:22-alpine AS runner

ENV NODE_ENV=production
ENV NITRO_HOST=0.0.0.0
ENV NITRO_PORT=3000
ENV NUXT_TELEMETRY_DISABLED=1

WORKDIR /app

COPY --link --from=builder /app/.output ./.output

EXPOSE 3000

CMD ["node", ".output/server/index.mjs"]
