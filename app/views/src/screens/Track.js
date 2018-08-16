import React from 'react'

class Track extends React.Component {
	name
	
	render() {
		return <li>{this.props.name}</li>
	}
}

export default Track;