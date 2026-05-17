# vk-photo-parser

A Go tool for downloading photos from VK (VKontakte) albums.

## Setup

1. Copy the example env file and fill in your credentials:

```bash
cp .env.example .env
```

2. Set the following environment variables in `.env`:

| Variable | Description |
|---|---|
| `VK_TOKEN` | Your VK API access token |
| `VK_OWNER_ID` | The owner ID of the album (user or community) |
| `VK_ALBUM_ID` | The album ID |

> **Note:** According to the VK API documentation, `VK_ALBUM_ID` must be prefixed with a `-` (minus sign) when the album belongs to a community/group. For example: `-123456789`.

## Getting a VK Token

You can obtain an access token via the [VK API documentation](https://dev.vk.com/api/access-token/getting-started).

## Running

```bash
go run ./cmd
```
