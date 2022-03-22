# Introduction to Attractify
In this tutorial, we want to show you how to integrate Attractify into your application.

For more detailed information about a single Attractify component, e.g. which parameters are needed to track an event and what effect they have, please take a look at our detailed [reference](https://github.com/inovex/attractify/blob/master/docs/reference.md).

## Creating a channel
In the first step, we will be creating a **channel**. This is used to assign users to their access origin and to better target them later. For example, if you have a website and an app, create a separate **channel** for each of them.

You can find the **channels** in the left menu of the Attractify backend. In our example, the channel `Web` is created for a website.

- `Name`: identify the **channel** later in Attractify.
- `Key`: reference the **channel** from the SDK and in API calls. We suggest a lowercase key without spaces and special chars, e.g. `web`.

![](/docs/assets/channel.gif)

## Creating a context
In order to capture, which boundary conditions the user brings along, you can define allowed context properties. For example, the country from which the user connects or the type of browser can be determined.

You can find the contexts in the left menu of the Attractify backend. In our example, the context is defined for the channel `Web`.

- `Key`: reference the **context property** from the SDK and in API calls. We suggest a lowercase key without spaces and special chars, e.g. `conn`.
- `Type`: defines the Datatype used for the property, e.g. `string`.
- `Regex Pattern`: used to ignore requests not matching the pattern.
- `Required`: If set, requests without that property will be ignored.

![](/docs/assets/context.gif)

## Generating an API key
In order for the Attractify script to authenticate itself, you need an API key. Your Attractify script won't work without an **API key**.

You can find the **API keys** in the left menu of the Attractify backend. You can create multiple API keys, for a given channel. In our example, an **API key** for the channel `Web` is created.

![](/docs/assets/api-key.gif)
## Creating an event
Now we create our first **event**. Events are used to track the user's behavior.

You can find the **events** in the left menu of the Attractify backend. In our example, the **event** tracks the visited page.

 - `Name`: reference the **event** from the SDK and in API calls, e.g. `pageName`.
 - `Description`: describe the use of the **event**.

#### Properties:
 - `key`:  reference the **event property** from the SDK and in API calls, e.g. `pageName`.
 - `Type`: datatype used for the property, e.g. `string`.
 - `Regex Pattern`: used to ignore requests not matching the pattern.
 - `Required`: If set, requests without that property will be ignored.

![](/docs/assets/events.gif)
## Tracking an event
In the example, the user visits our website. We **track** which page he visited and then look at the event log.

JavaScript code for the website tracking:

```
const apiConfig = {
  apiUrl: 'https://api.attractify.io/v1'
}
const attractify = new Attractify(
  'YOUR_AUTH_TOKEN_HERE',
  apiConfig
);

attractify.identify();
attractify.track('pageName', {
  pageName: window.location.href.replace(window.location.origin, ')
});
```
![](/docs/assets/tracking-event.gif)

## Creating a Computed Trait
As the user visits your website, Attractify can learn and provide a better and better user experience as time goes on. **Computed traits** play a crucial role in this. These are calculated dynamically and in real-time by incoming events.

You can find the **computed traits** in the left menu of the Attractify backend. In our example, the computed trait captures the last page change of the user and saves the timestamp.
 - `Name`:  identify the **computed trait** later in Attractify, e.g. `Last Activity`.
 - `key`:  reference the **computed trait** from the SDK and in API calls, e.g. `pageName`.
 - `Type`: defines how the value for the trait is computed, e.g. `Last Event`.
 - `Event`: the event which is used for the **computed trait**.
 - `Property Name`: event property which is computed, or use the timestamp.

 #### Event conditions:
 You can add event conditions to restrain the computation. For example, you could only save the timestamp if the pageName property in our example exists.
 - `Property`: property you want to check.
 - `Operator`: method to compare the values, e.g. `Exists`.
 - `Value`: value used for comparison.


![](/docs/assets/computed-trait.gif)

Now let's take a look at the result:
![](/docs/assets/computed-trait-example.gif)


## Creating a Custom Trait
You know your user better than Attractify - especially in the beginning. To store data from your systems directly in Attractify there are Custom Traits. They are filled by the SDK function "identify".

You can find the **custom traits** in the left menu of the Attractify backend. In our example, the custom trait captures the department an employee works.

 - `key`:  reference the **custom property** from the SDK and in API calls, e.g. `department`.
 - `Type`: datatype used for the property, e.g. `string`.
 - `Regex Pattern`: used to ignore requests not matching the pattern.
 - `Required`: If set, requests without that property will be ignored.

JavaScript code for attaching external data sources:

```
const traits = {
  department: 'marketing'
}
attractify.identify('user_id_123', 'user_id', traits);
```

This is how we create a custom trait:
![](/docs/assets/custom-trait.gif)

Now let's check if the Custom Trait is set when a new user comes to the page and is assigned to the Marketing department by another system:
![](/docs/assets/custom-trait-example.gif)

## Actions
**Actions** are played to actively improve the user experience. The preparation of data or texts is done in the frontend.

You can find the **actions** in the left menu of the Attractify backend. In the example, we will take a closer look at how, the user's most popular dress size can be set if a user chooses a product. The prerequisite is an event that tracks the dress size, as well as a computed trait that calculates the most common size.

 - `Name`:  identify the **action** later in Attractify, e.g. `Preselect Size`.
 - `type of action`:  reference the **action** from the SDK and in API calls, e.g. `size_preset`.
 - `Tag`: handle **actions** in your script by assigning tags.
 - `State`: define if the action is active or not.
 - `Properties`: send computed and custom trait properties with your **action**.
 - `Targeting`: who and when receives the **action**, e.g. `only user with department property set to marketing`.
 - `Capping`: how often should a user receive the action, e.g. `once per week`.
 - `Hooks`: define what happens if a user received the **action**.

JavaScript code for actions:

````
const getActions = async () => {
  const actions = await attractify.actions()
  for (a of actions) {
    if (a.type === 'size_preset') {
      adjustSize(a.properties.size); // preselect size in UI
    }
  }
};
````

![](/docs/assets/action.gif)

Now let's check if the action we just created works. For this we take a user who enters his size for the first time and when navigating to the store page again, the correct size is automatically suggested.
![](/docs/assets/action-example.gif)

## Analysis of Actions
After we have created events and custom traits and accumulated results, actions are played out. For example, it could be tracked that a user did not complete the payment process. The next time the user visits the website, a discount code is served to him via an action. To keep track of how well these actions are accepted, we offer a simple overview with the analyze tab.

Display the statistics for played out discount codes:

![](/docs/assets/analyze-action.gif)

## Process DSGVO requests
Data protection is more important than ever. At Attractify, we want to help you comprehensively provide user data to the appropriate user with minimal effort. You can easily delete or export his data. Another feature allows you to "lock" profiles. Once a profile has this status, no further data from Attractify will be assigned to that user.

Let's first look at how to delete any user data. Optionally, the data can be emailed to the user.
![](/docs/assets/gdpr-delete.gif)

The example shows how to lock a specific profile. In the future, Attractify will not store any further data about the user.
![](/docs/assets/gdpr-lock.gif)
