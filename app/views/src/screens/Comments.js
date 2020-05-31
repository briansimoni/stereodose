import React from 'react';

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faTrash } from '@fortawesome/free-solid-svg-icons';

class Comments extends React.Component {
  constructor(props) {
    super(props);

    this.textArea = React.createRef();

    this.state = {
      value: ''
    };
  }

  handleChange = event => {
    this.setState({ value: event.target.value });
  };

  render = () => {
    const { comments, onDeleteComment, user } = this.props;
    return (
      <div className="comments">
        <ul className="list-group">
          {comments.map(comment => {
            return (
              <li className="list-group-item" key={comment.ID}>
                <div>
                  <strong>{comment.displayName}</strong>
                  <br />
                  {comment.content}
                  {this.isUserComment(comment.userID) && (
                    <button
                      className="btn btn-danger delete-button"
                      onClick={() => {
                        onDeleteComment(comment.ID);
                      }}
                    >
                      <FontAwesomeIcon icon={faTrash}/>
                    </button>
                  )}
                </div>
              </li>
            );
          })}
        </ul>
        {user && (
          <div className="form-group">
            <label htmlFor="comment-textarea">Leave a Comment</label>
            <textarea
              ref={this.textArea}
              className="form-control"
              id="comment-textarea"
              rows="3"
              onChange={this.handleChange}
            />
            <button
              type="submit"
              className="btn btn-primary mb-2"
              onClick={() => {
                this.submitComment(this.state.value);
              }}
            >
              Submit
            </button>
          </div>
        )}
      </div>
    );
  };

  // submitComment wraps the parent function and clears the text after a button click
  submitComment = async text => {
    await this.props.onSubmitComment(text);
    this.textArea.current.value = '';
  };

  // isUser comment will return a boolean indicating whether
  // the userID attached to the comment matches the current User ID
  isUserComment = commentUserID => {
    const user = this.props.user;
    if (!user) {
      return;
    }
    if (user.ID === commentUserID) {
      return true;
    }
    return false;
  };
}

export default Comments;
