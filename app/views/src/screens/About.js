import React from 'react';

class About extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      feedbackSubmitted: false,
      requestInFlight: false,
      goodExperience: null,
      otherComments: '',
    }
  }

  handleChange = (event) => {
    // Using a dynamic property names to not have to create a switch statement
    // and write this line a whole bunch of times. The name is the HTML "name"
    // attribute on the corresponding HTML element.
    this.setState({ [event.target.name]: event.target.value });
  }

  handleSubmit = async (event) => {
    event.preventDefault();
    console.log('handling')
    if (this.state.requestInFlight || this.state.feedbackSubmitted) {
      console.log('returning');
      return;
    }
    this.setState({ requestInFlight: true });

    const feedback = this.state;
    feedback.goodExperience = feedback.goodExperience === "true" // switch from string to boolean

    try {
      await this.createFeedback();
      this.setState({
        requestInFlight: false,
        feedbackSubmitted: true,
      })
    } catch (err) {
      alert(err.message);
    }
  }

  createFeedback = async () => {
    const response = await fetch('/api/feedback/', {
      headers: {
        'Content-Type': 'application/json'
      },
      method: 'POST',
      body: JSON.stringify(this.state),
    })

    if (response.status !== 200) {
      const text = await response.text();
      throw new Error(`Problem creating feedback: ${text}`);
    }
  }

  render() {
    return (
      <div className="row">
        <div className="col">
          <h1 className="about-header">About</h1>
          <p className="about-body">
            Stereodose is a reincarnation of the web/mobile application that closed down back in 2016. It provided a way
            for people to discover music that mainstream services just don't offer. You can share your playlists from
            Spotify and let the community vote the best playlists to the top.
          <strong>
              {' '}
            Stereodose is currently in beta. It is open source so you can see progress and even contribute on
            <a href="https://github.com/briansimoni/stereodose"> GitHub</a>
            </strong>
          </p>

          <p>
            Please report issues to{' '}
            <a href="https://github.com/briansimoni/stereodose/issues">
              https://github.com/briansimoni/stereodose/issues
          </a>
          </p>

          <h1>Compatibility</h1>
          <p>
            You need to have Spotify Premium for the player to function. Additionally, the Spotify Web SDK only supports
            certain browsers. While it does seem to work okay on many mobile browsers, it isn't officially supported. See
          <a href="https://developer.spotify.com/documentation/web-playback-sdk/#supported-browsers">
              {' '}
            https://developer.spotify.com/documentation/web-playback-sdk/#supported-browsers
          </a>
          </p>
          <p>Stereodose is coming to the iOS app store in the next few months. Be sure to check back!</p>

          <h1>Send Feedback</h1>
          <form onSubmit={this.handleSubmit}>
            <h4>Do You Like Sterodose.app?</h4>
            <div className="form-check">
              <input onChange={this.handleChange} required="true" className="form-check-input" type="radio" name="goodExperience" id="yes-radio" value={true} />
              <label className="form-check-label" htmlFor="yes-radio">
                Yes
            </label>
            </div>
            <div className="form-check">
              <input onChange={this.handleChange} required="true" className="form-check-input" type="radio" name="goodExperience" id="no-radio" value={false} />
              <label className="form-check-label" htmlFor="no-radio">
                No
              </label>
            </div>


            <div className="form-group">
              <label htmlFor="other-comments">Any additional comments?</label>
              <textarea onChange={this.handleChange} name="otherComments" className="form-control" id="other-comments" rows="3"></textarea>
            </div>

            {/* Conditionally render what the success button looks like*/}
            {this.state.requestInFlight &&
              <button className="btn btn-primary">Sending Feedback</button>
            }

            {this.state.feedbackSubmitted &&
              <button type="submit" className="btn btn-success">Thanks For Submitting</button>
            }
            {!this.state.feedbackSubmitted && !this.state.requestInFlight &&
              <button type="submit" className="btn btn-primary">Submit</button>
            }

          </form>
        </div>
      </div>
    );
  }

}

export default About
