import React from 'react';
import { Fragment } from 'react';

class Feedback extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      formData: {
        goodExperience: null,
        otherComments: '',
      },
      characterLength: 0,
      feedbackSubmitted: false,
      requestInFlight: false,
    }
  }

  handleChange = (event) => {
    const state = this.state;
    state.formData[event.target.name] = event.target.value;
    this.setState(state);
  }

  handleSubmit = async (event) => {
    event.preventDefault();
    if (this.state.requestInFlight || this.state.feedbackSubmitted) {
      return;
    }
    this.setState({ requestInFlight: true });

    const feedback = this.state.formData;
    feedback.goodExperience = feedback.goodExperience === 'true' // switch from string to boolean

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
      body: JSON.stringify(this.state.formData),
    })

    if (response.status !== 200) {
      const text = await response.text();
      throw new Error(`Problem creating feedback: ${text}`);
    }
  }

  render() {

    if (this.state.feedbackSubmitted) {
      return (
        <div class="alert alert-success" role="alert">
          Thanks for submitting feedback!
        </div>
      )
    }
    return (
      <Fragment>
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
            <span id="feedback-char-limit">Limited to 10,000 characters</span>
            <textarea onChange={this.handleChange} name="otherComments" className="form-control" id="other-comments" rows="3" maxLength="10000"></textarea>
          </div>

          {/* Conditionally render what the success button looks like*/}
          {this.state.requestInFlight &&
            <button className="btn btn-primary disabled">Sending Feedback...</button>
          }

          {!this.state.feedbackSubmitted && !this.state.requestInFlight &&
            <button type="submit" className="btn btn-primary">Submit</button>
          }

        </form>
      </Fragment>
    )
  }
}

export default Feedback