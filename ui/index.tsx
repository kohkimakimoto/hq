import React from 'react'
import ReactDOM from 'react-dom'

declare var process: any;
if (process.env.NODE_ENV === 'development') {
  console.log('This is development mode!');
}

ReactDOM.render(<div>aaaaaaaaaaaa</div>, document.getElementById('app'));
