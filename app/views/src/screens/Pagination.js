import React from 'react';
import { Link } from 'react-router-dom';

export default function Pagination(props) {
  const params = new URLSearchParams(window.location.search);
  const page = params.get('page') !== null ? parseInt(params.get('page')) : 1;
  const { resultsPerPage, match, playlists } = props;

  return (
    <div className="row justify-content-center">
      <nav aria-label="Playlists search results page.">
        {page === 1 && ( // first page and previous button disabled
          <ul className="pagination">
            <li className="page-item disabled">
              <Link to="" tabIndex="-1" className="page-link">
                Previous
              </Link>
            </li>

            <li className="page-item active">
              <Link to={`${match.url}?page=${page}`} className="page-link">
                {page} <span className="sr-only">(current)</span>
              </Link>
            </li>

            {playlists.length > resultsPerPage && (
              <li className="page-item">
                <Link className="page-link" to={`${match.url}?page=${page + 1}`}>
                  {page + 1}
                </Link>
              </li>
            )}

            {playlists.length > resultsPerPage && (
              <li className="page-item">
                <Link className="page-link" to={`${match.url}?page=${page + 1}`}>
                  Next
                </Link>
              </li>
            )}
          </ul>
        )}

        {page > 1 &&
        playlists.length > resultsPerPage && ( // previous button + all the other buttons
            <ul className="pagination">
              <li className="page-item">
                <Link className="page-link" to={`${match.url}?page=${page - 1}`}>
                  Previous
                </Link>
              </li>

              <li className="page-item">
                <Link to={`${match.url}?page=${page - 1}`} className="page-link">
                  {page - 1}
                </Link>
              </li>

              <li className="page-item active">
                <Link to={`${match.url}?page=${page}`} className="page-link">
                  {page} <span className="sr-only">(current)</span>
                </Link>
              </li>

              <li className="page-item">
                <Link to={`${match.url}?page=${page + 1}`} className="page-link">
                  {page + 1}
                </Link>
              </li>

              <li className="page-item">
                <Link className="page-link" to={`${match.url}?page=${page + 1}`}>
                  Next
                </Link>
              </li>
            </ul>
          )}

        {playlists.length < resultsPerPage &&
        page !== 1 && ( // we are on the last page
            <ul className="pagination">
              <li className="page-item">
                <Link className="page-link" to={`${match.url}?page=${page - 1}`}>
                  Previous
                </Link>
              </li>

              <li className="page-item">
                <Link to={`${match.url}?page=${page - 1}`} className="page-link">
                  {page - 1}
                </Link>
              </li>

              <li className="page-item active">
                <Link to={`${match.url}?page=${page}`} className="page-link">
                  {page} <span className="sr-only">(current)</span>
                </Link>
              </li>

              <li className="page-item disabled">
                <Link to="" tabIndex="-1" className="page-link">
                  Next
                </Link>
              </li>
            </ul>
          )}
      </nav>
    </div>
  );
}
