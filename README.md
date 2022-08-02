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

## 🔎 About Attractify
The amount of information a user is confronted with during their search for a specific thing is sometimes overwhelming. That's why we think it is essential to personalize the web and in app experience for your users during their journey.

There are different approaches to solve these challenges. You can integrate web tracking, analyze the data in real time and try to predict the user's journey. However, you still need a service that takes over the evaluation and another service that then personalizes the experience for the user on your web site or app.

Yes, this is possible, but we see two problems here:

- The systems need to be extremely well connected.
- In times of GDPR and CCPA, such sensitive data should not reside with a third-party provider.

And these are the reasons why we developed Attractify. We needed a system that would allow us to personalize websites and apps without having to put the data in someone else's hands.


# 🚀 How to get started
The Attractify platform consists of two components. The API for developers and the UI for marketing manager.

We provide a Docker Compose file that includes all the bits and pieces to get you started.

```
curl -o compose.yaml "https://raw.githubusercontent.com/inovex/attractify/master/docker-compose.demo.yml" ; mkdir -p server/testdata/fixtures ; curl -o server/testdata/fixtures/postgres.sql "https://raw.githubusercontent.com/inovex/attractify/master/server/testdata/fixtures/postgres.sql"  ; docker-compose up
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

# ❓ Documentation
If you want to testdrive Attractify or deploy it in production, please see our [demo and installation guide](/docs/installation.md). To understand the various attractify features you can have a look at our [getting started guide](/docs/getting-started.md).
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

