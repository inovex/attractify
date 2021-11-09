let baseUrl = ''
if (process.env.NODE_ENV === 'production') {
  baseUrl = 'https://api.attractify.io/v1'
} else {
  baseUrl = 'http://127.0.0.1:8080/v1'
}

import Fifo from 'localstorage-fifo'

class Attractify {
  constructor(authToken) {
    this.state = { isIdentified: false }
    this.authToken = authToken
    this.context = null
    this.queue = new Fifo({ namespace: 'trackings' })

    this.loadState()

    if (!this.state.userId) {
      this.state.userId = this.uuidv4()
      this.state.isAnonymous = true
      this.saveState()
    }

    window.setInterval(() => {
      this.queueWorker()
    }, 500)
  }

  queueWorker() {
    const tasks = this.queue.get()
    for (const [k, v] of Object.entries(tasks)) {
      this.request('POST', '/track', v).then(() => {
        this.queue.remove(k)
      })
    }
  }

  setContext(arg) {
    this.context = arg
  }

  loadState() {
    try {
      this.state = JSON.parse(localStorage['_a8y'])
    } catch (_) {
      this.state = {}
    }
  }

  saveState() {
    localStorage['_a8y'] = JSON.stringify(this.state)
  }

  uuidv4() {
    return ([1e7] + -1e3 + -4e3 + -8e3 + -1e11).replace(/[018]/g, (c) =>
      (
        c ^
        (crypto.getRandomValues(new Uint8Array(1))[0] & (15 >> (c / 4)))
      ).toString(16)
    )
  }

  // Wrap fetch
  request(method, path, params) {
    const request = {
      method: method,
      headers: new Headers({
        Authorization: 'Bearer ' + this.authToken,
      }),
    }

    if (params) {
      request.body = JSON.stringify(params)
      request.headers.append('Content-Type', 'application/json')
    }

    return fetch(`${baseUrl}${path}`, request)
  }

  identify(userId, type = 'user_id', traits = null) {
    if (!userId || userId.length === 0) {
      if (!this.state.isAnonymous) {
        this.state.previousUserId = null
        this.state.isAnonymous = true
        this.state.userId = this.uuidv4()
      }
    }

    if (userId && userId.length > 0) {
      this.state.previousUserId = this.state.userId
      this.state.userId = userId
      this.state.isAnonymous = false
    }

    const params = {
      userId: this.state.userId,
      type: type,
      isAnonymous: this.state.isAnonymous,
      traits: traits,
    }

    if (this.state.userID !== this.state.previousUserId) {
      params.previousUserId = this.state.previousUserId
    }

    const res = this.request('POST', '/identify', params)
    this.state.isIdentified = true
    this.state.previousUserId = null

    this.saveState()
    return res
  }

  getContext() {
    if (typeof this.context === 'function') {
      return this.context()
    } else if (typeof this.context === 'object') {
      return this.context
    }

    return null
  }

  // Get actions from API
  loadActions(filters = {}) {
    return new Promise((resolve, reject) => {
      const params = {
        userId: this.state.userId,
        context: this.getContext(),
      }
      if (filters.tags) {
        params.tags = filters.tags
      }
      if (filters.type) {
        params.type = filters.type
      }

      this.request('POST', '/actions', params)
        .then((res) => {
          res.json().then(resolve).catch(reject)
        })
        .catch(reject)
    })
  }

  // Returns actions that match all filters
  actions() {
    let filters = {}
    let args = Array.from(arguments)
    if (args.length === 1 && typeof args[0] === 'object') {
      filters = args[0]
    } else {
      filters.tags = args
    }

    return new Promise((resolve, reject) => {
      this.loadActions(filters)
        .then((res) => {
          const actions = res.map((a) => {
            return this.prepareAction(a)
          })

          resolve(actions)
        })
        .catch(reject)
    })
  }

  async track(event, properties) {
    const params = {
      event: event,
      userId: this.state.userId,
      properties: properties,
      context: this.getContext(),
    }

    this.queue.set(this.uuidv4(), params)
  }

  // Track event
  act(event, actionId, properties) {
    const params = {
      event: event,
      actionId: actionId,
      userId: this.state.userId,
      properties: properties,
    }

    return this.request('POST', '/actions/act', params)
  }

  // Create callbacks and prepare slots
  prepareAction(action) {
    const me = this

    const actCb = (name, actionId, args) => {
      return new Promise((resolve, reject) => {
        if (args.length === 0) {
          me.act(name, actionId)
            .then((res) => res.json())
            .then(resolve)
            .catch(reject)
          return
        } else if (typeof args[0] === 'function') {
          // Argument is a callback
          const callRes = args[0]()
          if (callRes instanceof Promise) {
            callRes.then((res) => {
              me.act(name, actionId, res)
                .then((res) => res.json())
                .then(resolve)
                .catch(reject)
            })
          } else {
            me.act(name, actionId, callRes)
              .then((res) => res.json())
              .then(resolve)
              .catch(reject)
          }
          // Argument is not a callback
        } else {
          if (args.length > 0) {
            me.act(name, actionId, args[0])
              .then((res) => res.json())
              .then(resolve)
              .catch(reject)
            return
          }
        }
      })
    }

    return {
      id: action.id,
      type: action.type,
      callbacks: {
        show() {
          return new Promise((resolve, reject) => {
            me.act('show', action.id).then(resolve).catch(reject)
          })
        },
        hide() {
          return new Promise((resolve, reject) => {
            me.act('hide', action.id).then(resolve).catch(reject)
          })
        },
        decline() {
          return actCb('decline', action.id, arguments)
        },
        accept() {
          return actCb('accept', action.id, arguments)
        },
      },
      properties: action.properties,
    }
  }
}

window.Attractify = Attractify
