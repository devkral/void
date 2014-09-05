HOST= "http://10.40.0.19"

window.App = Ember.Application.create
    LOG_TRANSITIONS: true
    LOG_TRANSITIONS_INTERNAL: true
    authstring: ""
    sessionUser : ->
        $.base64.atob(@authstring).split(":",1)[0]

class App.ApplicationStore extends DS.Store

DS.RESTAdapter.reopen
    namespace: 'rest'
    host: HOST
    +volatile
    headers : ->
        Authorization : "Basic "+App.authstring

mixin CapitalAttrs
    keyForAttribute: (attr) -> attr.charAt(0).toUpperCase()+attr.slice(1)

App.Router.map ->
    @resource "buildings", ->
        @route "list"
        @route "map"
    @resource "building", {path: '/building/:id'}, ->
        @route "edit"
        @route "view"
        @resource "comments"
    @resource "user", {path: '/user/:id'}, ->
        @route "edit"

class App.IndexRoute extends Ember.Route
    model: -> Ember.Object.create
                  login: Ember.Object.create
                      username: ""
                      password: ""
                  building: Ember.Object.create
                      name: ""
                      address: ""

class App.IndexController extends Ember.ObjectController
    actions:
        post_building: ->
            Ember.Logger.debug "post"
        login : ->
            Ember.Logger.debug "login"

class App.ApplicationRoute extends Ember.Route
    model : -> null

class App.ApplicationController extends Ember.Controller
    actions:
        createBuilding: ->
            Ember.Logger.debug "create"

