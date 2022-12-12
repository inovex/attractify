export default {
  channels: {
    title: 'What are channels?',
    body: `
Channels define the source of your users. If you have a website
and an app, you should create two separate channels. This makes
it easier to examine your users behaviour on a per channel basis
and helps you target each channel with different marketing
actions. If a single user is identified by a consistent ID on
your website and in your app, Attractify will merge all events
of this user over both channels.

# Parameters

* \`Name\` - The name is only used to describe the channel.
* \`Key\` - The key is used internally to track the source of events and
destinations of marketing actions. We suggest to use a single,
descriptive word like *web*, *app*, *email*.

It is not possible, to change the key after the channel has been created.

# Actions

* Edit channel
* Delete channel

# Deleting channels

If you delete a channel without stopping the event tracking or
changing your audiences or marketing actions to a new channel,
they will stop working. The best way is to leave a channel
unused instead of deleting it.`
  },
  events: {
    title: 'What are events?',
    body: `
Events are the essential base for learning more about your user's behaviour.
Each important action on your website or in your app should be reflected by an
event definition.

Events can contain additional properties with structured data to provide more
context and to better understand your user.

Please note, that if you delete an event, we still keep all tracked instances of it.

# Actions

* Edit event
* Delete event

`
  },
  invalidevents: {
    title: 'What are invalid events?',
    body: `
In the invalid event section you can find all events where the context or properties
do not match the defined event.

# Actions

* View details
* Delete event

`
  },
  event: {
    title: 'Defining an event',
    body: `
An event has a name, a description and a list of properties. The
name should follow a consistent naming convention to prevent
confusion. We suggest something like \`Subject Action\` e.g. \`Page Viewed\`,
\`Product Added To Cart\` or \`Checkout Completed\`.

You could also use \`camelCase\` like \`pageViewed\`, \`productAddedToCart\` or \`checkoutCompleted\`.

# Properties

Properties are values that could be attached to an event to
track additional information. When a user adds a product to
the shopping cart, you can add the product's ID and price as a
property. The more detailled the information you have tracked
is, the better you can fine tune your audiences or marketing
actions. You need to defined the properties ahead of tracking so
that Attractify can validate the incoming properties against the
definition. This ensures that only valid information is stored.

# Parameters

* \`Name\` - The name of the event. Must be unique. Please use a consistent naming schema.
* \`Description\` - A description to better understand, what the event does.

# Changing events

After an event is  in use, you should be careful about changing its name or properties.
Adding additional properties is no problem, but before removing or modifying a
property, please make sure it is not used by computed traits or actions.

Changing an event name that is used by computed traits or audiences should be avoided.
`
  },
  properties: {
    title: 'Property parameters',
    body: `
A single property needs to be identified by a unique key. We cast the value to the given data type.

* \`Key\` - The name of the property. Please use camelCase as naming schema.
* \`Type\` - A data type for the property.
* \`Regex Pattern\` - If you want to validate the value before casting it to the given data type, you could add you own regex.
* \`Required\` - Indicates if the property must be supplied and valid.
`
  },
  contexts: {
    title: 'What are contexts?',
    body: `
Contexts can contain additional information for events and actions. Maybe you want to target users from different countries or with different browsers? Or you want to know on which page on your website a specific event has occured. Then you can supply this additional information as context.

Each request that contains a context will be validated against the context definition.
`
  },
  context: {
    title: 'How to define a context?',
    body: `
Every context is associated with a channel. You can only assign one context to each channel as we're resolving the context
definition on every event or action call based on its channel.
`
  },
  eventLog: {
    title: 'What is the event log?',
    body: `
The event log contains all events that have been tracked over all
channels. You can filter these by event name or search for events
from a specific user.

# Actions

* Show event with properties and context
* Delete event
`
  },
  profiles: {
    title: 'What are user profiles?',
    body: `
A user profile combines all events, identities and traits under a single persona.
User profiles are automatically created if a tracking event or action request is
received.

Each user profile will be assigned a unique ID to identify this user. If traits
for a user were updated, the last updated timestamp is refreshed.

If you delete a profile, all associated identities and events will be deleted too.

# Actions

* Show profile
* Refresh computed traits
* Delete profile
`
  },
  profile: {
    title: 'What is a user profile?',
    body: `
A user profile contains all custom and computed traits as well as all identities
and events for a single user.

## Identities

A user profile can have multiple identities from one or more channels. Each identity
has its own ID, a type, the channel it originates from and a flag if it is anonymous or not.
The type can be freely chosen trough the \`identify\` API call and specifies what kind of
ID is used. Examples are

* \`email\` - for an email address
* \`userId\` - if you use your internal user ID
* \`mobileNumber\` - for a mobile subscriber number

and many more.

## Custom traits

Custom traits are pieces of information that are already known to the organization
and can be supplied using the \`identify\` API.

## Computed traits

Compared to custom traits, computed traits are dynamically refreshed when a new
tracking event is received.

# Actions

* Refresh profile
`
  },
  customTraits: {
    title: 'What are custom traits?',
    body: `
Sometime your organization knows a lot more about a user than Attractify does,
because the user is already a customer. In this case custom traits can help to,
share that knowledge with Attractify.

You can supply these information through the \`identify\` API call. But to make
sure, that only valid data is imported and stored in attractify, you need to
create a schema. This ensures that only predefined properties with the corresponding
data type are accepted.

The specifics of those properties are the same as in the event or
context definition.

Changing custom traits after user profiles have been created does not pose
a problem. Existing custom traits of a user profiles are not modified and new user
profiles will get the new properties.
`
  },
  computedTraits: {
    title: 'What are computed traits?',
    body: `
In comparison to static custom traits, computed traits are dynamically computed based
on incoming tracking events. Computed traits represent information that are
collected in realtime and may change over time.

They are perfectly suited to better understand users and can be used
to directly target actions or build audiences from.

If a computed trait is deleted, the user profiles will still contain that
trait, but will no longer be updated.

# Actions

* Edit computed trait
* Delete computed trait
`
  },
  computedTrait: {
    title: 'Defining a computed trait',
    body: `
A computed trait defines a single property of a user profile that is updated
as soon as an event is tracked.

# Parameters

* \`Name\` - The name of the computed property.
* \`Key\` - A key that is used to access this property later. We recommend \`camlCamse\` as a naming schema.
* \`Type\` - The method that is used to compute the trait.
* \`Event\` - The event that should be used for computation.
* \`Property Name\` - Some methods require a single property that is selected here.
* \`Aggregation Type\` - For aggregations, you need to specify a type.
* \`Use timestamp instead of property value\` - For selecting a single event, you can either return the timestamp of the event or its value.
* \`Property Name\` - If you want to return the property's value instead of its timestamp, select the property here.
* \`Event conditions\` - Events that are used for computation can meet certain conditions.
`
  },
  audiences: {
    title: 'What are audiences?',
    body: `
Audiences are a way to group user profiles together based on certain criteria.
Sometimes you have tough requirements for a specific target group in your actions.
Then you have to use audiences instead of targeting users by their traits.

Audiences can have one of two types:

* Events
* Funnel

Events and Funnel are configured the same way. The only difference is, that in a funnel
the order of the events must match the definition whereas in a list of events, the order
does not play an important role.

If you want to make sure, that a user has a specific journey, use \`funnels\`, otherwise
\`events\`.


# Actions

* Refresh audience
* Edit audience
* Delete audience
`
  },
  audience: {
    title: 'Defining an audience',
    body: `
An audience is a dynamic group of user profiles that is computed based on the
given criteria. Once an audience has been saved, its type cannot be changed.

# Parameters

* \`Type\` - Does the order of events matter, than choose \`Funnel\` else \`Events\`.
* \`Name\` - A name for the audience.
* \`Description\` - We suggest to describe your audience, so that others understand the idea behind it.
* \`Include anonymouse profiles\` - You can in- or exclude profiles, that have not received an \`identify\` call.
* \`Audience definition\` - Add one or more events with conditions, custom or computed traits.
`
  },
  actions: {
    title: 'What are actions?',
    body: `
An action is an instruction pushed to the end user's device to

* show a specific piece of content
* call for an interaction
* modify the user's UI
* ask for information
* and many more.

Actions are building blocks to optimize the overall user experience.
They only contain the content and meta information but no UI elements or styling.
It is up to the frontend to define the appearance of an action.

But actions are not pushed to all users but instead are carefully targeted
based on the given definitions.

An action implies a reaction. These are events which are triggered on the users
device and help to measure if an action was successful or not. They are also
used with the capping feature to prevent an action to be shown multiple times
to the same user if he has declined it.

# Actions

* Analyze specific action
* Set action to active, inactive or staging
* Edit action
* Delete action
* Duplicate action
`
  },
  action: {
    title: 'Defining an action',
    body: `
An action has metadata and detailled configuration options as shown below.

# Parameters

* \`Name\` - A name for your action.
* \`Type\` - The type is used to signal the client that receives an action,
how to display or handle the action properly. This should be in \`camelCase\`.
* \`Tags\` - Tags are used to offer the client a possibility to request
actions based on specific criteria. Tags can also be used to test different
variants of an action.
* \`State\` - Specifies if the action is *active*, *inactive* or in a *staging*
state, which is used for testing.

An action is composed of different parts which are explained below:

# Properties

Properties can be used to supply content or configuration to the corresponding
UI element, which renders the action on the user's device. You can see them as parameters,
which are customizable and/or can contain placeholders.

Different properties can also be shown on different channels. This is very useful
to display an action differently on various devices. Maybe your text should be
shorter on mobile and more explanative on desktop. Just specify
the same property name with different values for different channels.

Properties can either be static text, or you can access custom or computed traits
from the current user.

# Targeting

To select which user should get which action, you can use targeting. The different
targeting options work like filters put over each other. The more targeting options
you sepcify, the more selective your targeting is.

The simplest form of targeting is to select one or more channels, your actions
should be pushed to. To further restrict the number of affected users you can
select an existing audience that contains a list of users.

If your action should only be active in a certain time range, start or end on a
specific date, you can configure this using the time range option.

Last but not least you can filter users by their custom oder computed traits
or the context, that is provided to the request.

# Capping

Instead of spamming users with the same action over and over again, you can
set a specific limit on the number of *shows*, *hides*, *declines* and *accepts*.
You can also specify a time window to consider only the last X days.

Capping can also be used to define a number of limits not on a per user basis,
but for all users in total. This is very useful to offer a specific number of
coupons and not exceed that number.

You can also specify different quotas for different channels. This is especially
helpful if you want to separte web and mobile as web users are more likely
to ignore some information dialogs compared to mobile users.

# Hooks

Hooks offer the possibility to attach external webhooks to actions. These
are called once the selected reaction has been received from the user. An example
would be, when the user accepts a coupon, the device will send the \`accept\` reaction
to Attractify and Attractify will call the configured webhook which can redeem the coupon.

You can also specify to create an event once a reaction has been received. This is useful
if you want to add users later which have reacted to a specific action to an audience.

# Testusers

To be able to test an action in production, you can add a test user for this action.
When an action is set to *staging*, you will still be able to test that action with
the configured test user. Other users won't get that action pushed to their devices.

If you enable *skip targeting*, the action will be testable, regardless of the
filters that are configured under targeting.
`
  },
  actionTypes: {
    title: 'What are action types?',
    body: `
An action template is a predefined type of action which defines which properties the action has.

# Actions
* Edit action template
* Delete action template
`
  },
  actionType: {
    title: 'Defining an action type',
    body: `

    Properties can be used to supply content or configuration to the corresponding
    UI element, which renders the action on the user's device. You can see them as parameters,
    which are customizable and/or can contain placeholders.
    
    Different properties can also be shown on different channels. This is very useful
    to display an action differently on various devices. Maybe your text should be
    shorter on mobile and more explanative on desktop. Just specify
    the same property name with different values for different channels.
    
    Properties can either be static text, or you can access custom or computed traits
    from the current user.

    Whenever an action type is in use. A new version is created when you edit a feature.
`
  },
  reactions: {
    title: 'What are reactions?',
    body: `
Reactions are events, which are triggerd after a user has interacted with an action.
This is used to measure the acceptance and popularity of an action and helps to
further optimize that action.

# Types of reactions

* \`Delivered\` - Attractify has delivered an action to the user.
* \`Shown\` - The action has been shown on the users device.
* \`Hidden\` - The user has hidden the action.
* \`Declined\` - The user has declined the action (stronger than hidden).
* \`Accepted\` - The user has accepted the action.

# Actions

* Show reaction
* Delete reaction
`
  },
  analyze: {
    title: 'How to analyze an action?',
    body: `
To find out, if an action performed as anticipated, you can use the analyze feature.
After selecting an action and optionally a date range and resulution, you see a
drilldown of the reactions.

# Events

Events is a simple summary how often a single event type has been tracked for
the given action. As events drip from top to bottom, you see, at which step
the action stopped performing well.

# Rates

Rates expresses a ratio between the last and the current event. Maybe your action
has been shown 10 times but is only accepted 5 times. Then you have a ratio of 2:1 and a
conversion drop of 50 percent. This could be an indicator to optimize your action.

# Reach

Reach is a per channel summary which event has been tracked how often. It helps
you to see, which is the dominant channel for that action.
`
  },
  api: {
    title: 'Configuring API access',
    body: `
If you want to interact with Attractify you have multiple possibilities to do so.

# Auth tokens for client integrations

To enable access to your clients (website, app, ...) you need to create an auth
token. Each channel must use a spearate auth token. This is necessary to assign
incoming requests to the corresponding channel.

# Webhook signature key

Outgoing webhook requests for actions are signed. To verify the signature you
need to obtain the webhook signature key, encoded as hex.
The signature is an \`HMACSHA256(body)\` over the request's body with the given key.

# Platform API access

To access APIs of the Attractify platform, which is the admin interface you are
currently viewing, you need a JWT token. If you press *Generate token* we will
generate a new token for you to use in your API client.
`
  },
  users: {
    title: 'Manage users',
    body: `
You can invite additional users or delete current ones. If you invite users, they'll
receive an invitation email and can then sign up.

# Parameters

* \`Email\` - The user's email that should be invited.
* \`Role\` - The role the user should have.

# Roles

* \`Admin\` - Has access to all features.
* \`Marketeer\` - Can access everything except user management and organizational settings.
`
  },
  organization: {
    title: 'Manage an organization',
    body: `
Here you can modify your organization's settings. Adjusting the timezone has an
impact on incoming events, as they are stored in the selected timezone.

# Parameters

* \`Name\` - Your organization's name.
* \`Timezone\` - The timezone your organization is operating in.
This is used to store trackings under this timezone.
`
  },
  privacy: {
    title: 'Privacy features',
    body: `
The Attractify privacy features help you to comply with the GDPR or DSGVO.

You have the following options:

* Export all stored data for a user.
* Delete all stored data for a user.
* Lock a user's identity to prevent further tracking.

# Parameters

* \`User ID\` - If you want to lock a user, simply enter the user's ID and further trackings for this ID will be blocked.
`
  },
  user: {
    title: 'Your personal settings',
    body: `
Here you can either update your name or email address or change your password. If
you want to change your password, make sure it is at least 8 chars long. We recommend
a longer password (> 20 chars) with a combination of numbers (upper- and lowercase), letters and special chars.
`
  },
  actionsimulation: {
    title: 'Action Simulation',
    body: `
  Here you can debug an action with an user.
  That way you can figure out what is causing an action not to be displayed.
  `
  }
}

