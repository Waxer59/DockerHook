<!-- TOC -->
* [DockerHook](#dockerhook)
  * [Get started](#get-started)
  * [Environment Variables](#environment-variables)
    * [`CONFIG_PATH`](#config_path)
    * [`PORT`](#port)
  * [WebHook](#webhook)
    * [Url](#url)
    * [Query Parameters](#query-parameters)
      * [`action`](#action)
      * [`token`](#token)
  * [Config Properties](#config-properties)
    * [`config`](#config)
      * [`labelBased`](#labelbased)
      * [`defaultAction`](#defaultaction)
    * [`auth`](#auth)
      * [`enable`](#enable)
      * [`tokensFile`](#tokensfile)
      * [`tokens`](#tokens)
  * [Examples](#examples)
<!-- TOC -->

# DockerHook

<p align="center">
  <img src="./docs/imgs/logo.jpg" width="512px" alt="DockerHook logo"/>
</p>

DockerHook is a way to manage your Docker processes using a webhook!

## Get started

```yaml

```

## Environment Variables

DockerHook uses the environment variables to allow you to customize the parameters

### `CONFIG_PATH`

* **Type**: `string`
* **Default**: `/etc/dockerhook/dockerhook.yml`

Specify the path where your configuration file is located

### `PORT`

* **Type**: `int`
* **Default**: `8080`

Specify the port that DockerHook will use.

## WebHook

### Url

`http://exampleurl.com/<docker-service-name>?<query-parameters>`

### Query Parameters

#### `action`

* **Type**: `'start' | 'stop' | 'restart' | 'pull'`
* **Default**: `'pull'`

Specify the action to be done in the service

#### `token`

* **Type**: `string`

Specify the access token to control Webhook execution.

Specifies the access token to control the execution of the Webhook if the [`enable`](#enable) property of [`auth`](#auth) is set to `true`.

## Config Properties

### `config`

#### `labelBased`

* **Type**: `boolean`
* **Default**: `false`

TODO

#### `defaultAction`

* **Type**: `'start' | 'stop' | 'restart' | 'pull'`
* **Default**: `'pull'`

### `auth`

#### `enable`

* **Type**: `boolean`
* **Default** `false`

#### `tokensFile`

* **Type**: `string`

#### `tokens`

* **Type**: `string`

## Examples

Check the following examples:

* [Example configuration](./docs/examples/exampleConfig.yml)
* [Example docker compose](./docs/examples/docker-compose.yml)
