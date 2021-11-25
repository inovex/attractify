![Attractify header](/assets/header.png)

We are developers and we hate to integrate marketing tools into websites and apps. We want clean APIs and no tools that generate garbage HTML that we have to bend into direction.
This is why we have built Attractify: an open-source customer data platform (CDP) like Segment but with an
integrated marketing automation engine (MAE). Attractify does all the heavy lifting of a CDP and MAE
while providing a simple API as well as SDKs that you can integrate into your website and/or app. In short, we use the headless approach every dev loves :)
For marketers we offer a comfortable UI to create and optimize their audiences and marketing automations.

And as we believe in data privacy, you can run Attractify for free on your own infrastructure
or use our fully managed Attractify hosting in our german data center.

<a href="https://www.inovex.de/" target="_blank"><img src="assets/sponsor.png"/></a>

The work on Attractify is sponsored by [inovex.de](https://inovex.de)

# 🧐 Wait, what does Attractify do?

First we help you to understand how your users behave across multiple channels. Then you can run your own cross-channel actions on them, like:

- Display notifications across channels and track if a notification has been viewed by the user. Then hide that notification on all other channels if the user has viewed it.
- Know if users have items in their cart when they return to your shop and display them a coupon code to push checkouts.
- Detect which is the desired cloth size of your user and pre-select them on subsequent product detail pages.
- Run discount code campaigns that are limited to a maximum number of codes.
- Segment users based on their price sensitivity and show them relevant offers.

We think developers know best how they should implement specific marketing campaigns. That is why Attracify does not offer predefined templates for these use cases. Instead we provide an API with SDKs that you can use to query if the current user is eligible for an action. Then you can run this action on the user and measure his interaction with it. This gives you the maximum flexibility and works on the web as well as in native apps.

# 🚀 How to get started

The Attractify platform consists of two components. The API for developers and the UI for marketing manager.

We provide a Docker Compose file that includes all the bits and pieces to get you started.

```
curl https://raw.githubusercontent.com/inovex/attractify/master/docker-compose.yml | docker-compose  -f /dev/stdin up
```

Or you can sign up for a free trial at [attractify.io](https://attractify.io).

Once Attractify is started, you can visit the usecase-shop under [http://localhost:8000](http://localhost:8000)

The API is available under [http://localhost:8080](http://localhost:8080). You can use the following credentials to login:

User: `demo@example.com`\
Password: `demo4321`

# 🧪 The Attractify Experimentation Loop

Our Attractify platform consists of four components that create the
Attractify Experimentation Loop. We call it that way as we think, that each marketing
campaign is an experiment and you need an easy way to adjust your campaign settings to
get the best results.

1. **Events:** To provide a great experience to your users, you need to get to know them first. Attractify offers an easy to use event tracking API with automatic data validation to understand your user's behaviour across channels and devices.
2. **Personas & Audiences**: You can group user profiles into personas and enrich them with custom and computed traits. Attractify can automatically segment your personas into audiences based on your own criteria.
3. **Actions**: You can do whatever you think improves the happiness of your users. Attractify helps you to coordinate which user should get what marketing action. Subsequently it measures their reaction to a campaign. Just query the Attractify API and it responds with the appropriate action for the current user.
4. **Insights**: Marketers love to see how their actions have performed. That is why insights helps to optimize their marketing efforts over time. Experiment, learn, repeat.

![inovex logo](/assets/experimentation-loop.jpg)
# 🕵️‍♀️ Data Privacy

You should own your user's data. That's why you can fully self host Attractify for free. And as we think GDPR could benefit us all, we have an integrated GDPR module to answer your users requests for data export, deletion and locking.

# 📃 License

Attractify and its SDKs are licensed under the MIT license which can be found in the [LICENSE](/LICENSE) file. We reserve the right to offer paid extensions later, that require the purchase of a commercial license. But the base platform will be free. Forever.

# 📌 Contributing

We love contributions. If you want to get started but don't know where exactly, <a href="mailto:marc.boeker+attractify@inovex.de?subject=Contributing to Attractify">send us a message</a> and we can give you a quick tour and talk about possible features/improvements we would love to see in Attractify.

# 💬 Feedback
We would love to hear your feedback on Attractify. Drop us a message in our `#support` channel on [Slack](https://attractify.slack.com) or [e-mail](mailto:info@attractify.io) us if you ran your first experiments.
# 💻 Stack

The backend is written in Go and the frontend in Vuejs. To store and process all the data, we use the following combination of awesome technology:

- [CockroachDB](https://www.cockroachlabs.com/) (for metadata, user profiles and identities)
- [ClickhouseDB](https://clickhouse.tech/) (for segmenting users into audiences)
- [Zookeeper](https://zookeeper.apache.org/) (for coordinating Clickhouse and Kafka replication)
- [Kafka](https://kafka.apache.org/) (for handling incoming tracking events)

