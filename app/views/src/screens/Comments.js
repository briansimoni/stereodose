import React from "react";
import Octicon from "react-octicon";


class Comments extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      value: ""
    };
  }

  handleChange = (event) => {
    this.setState({ value: event.target.value });
  }

  render = () => {
    const { comments, onSubmitComment, onDeleteComment } = this.props;
    return (
      <div className="comments">
        <ul className="list-group">
          {comments.map((comment) => {
            return (
              <li className="list-group-item" key={comment.ID}>
                <div>
                  <strong>{comment.displayName}</strong>
                  <br />
                  {comment.content}
                  {this.isUserComment(comment.userID) &&
                    <button className="btn btn-danger delete-button" onClick={() => { onDeleteComment(comment.ID) }}>
                      <Octicon name="trashcan" />
                    </button>
                  }
                </div>
              </li>
            )
          })}
        </ul>
        <div className="comment-text-area form-group">
          <label htmlFor="comment-textarea">Leave a Comment</label>
          <textarea className="form-control" id="comment-textarea" rows="3" onChange={this.handleChange}></textarea>
          <button type="submit" className="btn btn-primary mb-2" onClick={() => { onSubmitComment(this.state.value) }}>Submit</button>
        </div>
      </div>

    )
  }

  isUserComment = (commentUserID) => {
    const user = this.props.user;
    if (user.ID === commentUserID) {
      return true;
    }
    return false;
  }
}

export default Comments;