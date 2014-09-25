###########################################################
# © 2014 Daniel 'grindhold' Brendle and Contributors 
#
# This file is part of Void.
#
# Void is free software: you can redistribute it and/or
# modify it under the terms of the GNU Affero General Public License
# as published by the Free Software Foundation, either
# version 3 of the License, or (at your option) any later
# version.
#
# Void is distributed in the hope that it will be
# useful, but WITHOUT ANY WARRANTY; without even the implied
# warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR
# PURPOSE. See the GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public
# License along with Void.
# If not, see http://www.gnu.org/licenses/.
###########################################################

CLDR.defaultLocale = 'en'
Ember.FEATURES.I18N_TRANSLATE_HELPER_SPAN = false

window.App = Ember.Application.create
    LOG_TRANSITIONS: true
    LOG_TRANSITIONS_INTERNAL: true
    authstring: ""
    sessionUser : ->
        $.base64.atob(@authstring).split(":",1)[0]

class App.ApplicationStore extends DS.Store

DS.RESTAdapter.reopen
    namespace: 'rest'
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
                      captcha :""
    setupController : (c, m) ->
        @_super c, m
        c.loadCaptcha()

class App.IndexController extends Ember.ObjectController
    needs: ['application']
    captchaimg : ~> if @captchaid != "" then "/captcha/#{@captchaid}.png" else ""
    captchaid : ~> ""
    loadCaptcha: ->
        self = this
        $.ajax "/captcha/new",
            async:true
            dataType:"json"
            success: (data, status, xhr) ->
                self.captchaid = data.CaptchaId
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
                    captcha: @content.building.captcha
                    captchaid: @captchaid
            self = this
            b.save().then( ->
                Bootstrap.NM.push (Em.I18n.t 'index.postsuccess'), 'success'
                self.loadCaptcha()
            ).catch( (reason) ->
                if reason.status == 403
                  Bootstrap.NM.push (Em.I18n.t 'index.captchafail'), 'warning'
                else
                  Bootstrap.NM.push (Em.I18n.t 'index.postfailure'), 'danger'
            )
        login : ->
            App.authstring = $.base64.btoa @content.login.username+":"+@content.login.password
            self = this
            $.ajax "/auth",
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
                        Bootstrap.NM.push(Em.I18n.translations['index.loginerror'], 'warning')
            return


class App.ApplicationRoute extends Ember.Route
    beforeModel: (t) ->
        @controllerFor('application').loadAuthString()
    model : -> null

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
