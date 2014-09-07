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

class App.InvitationsNewRoute extends Ember.Route
    model : -> Ember.Object.create
                   email: ""
                   link: ""

class App.InvitationRoute extends Ember.Route
    model : (params) -> @store.find 'invitation', params.ident

class App.InvitationSerializer extends DS.RESTSerializer with CapitalAttrs
    keyForRelation: (key,name) -> key

class App.Invitation extends DS.Model
    email : DS.attr 'string'
    name : DS.attr 'string'
    organization : DS.attr 'string'
    password : DS.attr 'string'

class App.InvitationsNewController extends Ember.Controller
    inviteLink : ~> "#{HOST}/#/invitation/#{@content.link}"
    linkPresent : ~> @content.link != ""
    actions:
        invite : ->
            i = @store.createRecord 'invitation',
                  email: @model.email
            self = this
            i.save().then (x) ->
                self.content.link = x.id

class App.InvitationController extends Ember.ObjectController
    actions :
        redeem : ->
            @content.save()
            @transitionTo "index"
