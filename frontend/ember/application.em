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
    @resource "building", {path: '/building/:ident'}, ->
        @route "edit"
        @route "view"
        @resource "comments"
    @resource "user", {path: '/user/:id'}, ->
        @route "edit"
    @resource "invitations", ->
        @route "new"
    @resource "invitation", {path: '/invitation/:ident'}

class App.IndexRoute extends Ember.Route
    model: -> Ember.Object.create
                  login: Ember.Object.create
                      username: ""
                      password: ""
                  building: Ember.Object.create
                      street : ""
                      number : ""
                      city   : ""
                      zip    : ""
                      ownername : ""
                      ownerphone : ""
                      owneremail  : ""
                      area : 0
                      description : ""

class App.IndexController extends Ember.ObjectController
    needs: ['application']
    actions:
        post_building: ->
            b = @store.createRecord 'building',
                    street: @content.building.street
                    number: @content.building.number
                    city  : @content.building.city
                    zip   : @content.building.zip
                    ownername : @content.building.ownername
                    ownerphone : @content.building.ownerphone
                    owneremail : @content.building.owneremail
                    area        : @content.building.area
                    description : @content.building.descritpion
            b.save()
        login : ->
            App.authstring = $.base64.btoa @content.login.username+":"+@content.login.password
            self = this
            $.ajax HOST+"/auth",
                async: true
                dataType: "json"
                cache: false
                headers:
                    Authorization: "Basic "+App.authstring
                success: (data, status, xhr) ->
                    self.controllers.application.loggedin = data.Valid
                    if data.Valid
                        self.controllers.application.setAuthString App.authstring
                    else
                        self.controllers.application.resetAuthString()
            return


class App.ApplicationRoute extends Ember.Route
    model : -> null
    setupController : (controller, model) ->
        @_super controller,model
        @controllerFor('application').loadAuthString()

class App.ApplicationController extends Ember.Controller
    loggedin : ~> false
  
    setAuthString : (a) ->
        if sessionStorage
            sessionStorage.void_auth = a

    resetAuthString : ->
        if sessionStorage
            sessionStorage.removeItem "void_auth"

    loadAuthString : ->
        if sessionStorage and sessionStorage.void_auth
            App.authstring = sessionStorage.void_auth
            @loggedin = true

    actions:
        logout: ->
            @resetAuthString()
            App.authstring = ""
            @loggedin = false
            @transitionToRoute 'index'
            return

class App.FromNowView extends Ember.View
    tagName: 'time'
    template: Ember.Handlebars.compile '{{view.output}}'
    output: ~>
        return moment(@value).fromNow()
    didInsertElement: ->
        @tick()
    tick: ->
        f = ->
            @notifyPropertyChange 'output'
            @tick()
        nextTick = Ember.run.later(this, f, 1000)
        @set 'nextTick', nextTick
    willDestroyElement: ->
        nextTick = @nextTick
        Ember.run.cancel nextTick
Ember.Handlebars.helper 'fromNow', App.FromNowView
