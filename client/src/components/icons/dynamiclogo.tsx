import React from "react"

const CustomIcon = ({ width = 468, color = "black" }) => {
	// Calculate height based on the original aspect ratio
	const originalWidth = 468
	const originalHeight = 388
	const height = (width / originalWidth) * originalHeight

	return (
		<svg
			width={width}
			height={height}
			viewBox="0 0 468 388"
			fill="none"
			xmlns="http://www.w3.org/2000/svg"
		>
			<path
				d="M403.5 189H467.5C467.5 189 364 388 234 388C104 388 0 189 0 189H66.5C66.5 189 150.5 321.5 234 321.5C317.5 321.5 403.5 189 403.5 189Z"
				fill={color}
			/>
			{/* <path
				fillRule="evenodd"
				clipRule="evenodd"
				d="M235.5 288C286.586 288 328 246.586 328 195.5C328 144.414 286.586 103 235.5 103C184.414 103 143 144.414 143 195.5C143 246.586 184.414 288 235.5 288ZM236 275C280.183 275 316 239.183 316 195C316 150.817 280.183 115 236 115C191.817 115 156 150.817 156 195C156 239.183 191.817 275 236 275Z"
				fill={color}
			/> */}
			<path
				d="M168.885 199.662L233.191 135.79C234.747 134.245 237.257 134.241 238.819 135.781L303.559 199.653C306.105 202.164 304.327 206.5 300.75 206.5H171.703C168.133 206.5 166.351 202.178 168.885 199.662Z"
				fill={color}
			/>
			<rect
				x="210.5"
				y="195"
				width="50"
				height="55"
				rx="4"
				fill={color}
			/>
			<rect
				x="66"
				y="73"
				width="339"
				height="25"
				rx="12.5"
				fill={color}
			/>
			<rect x="103" y="44" width="264" height="20" rx="10" fill={color} />
			<rect
				x="141"
				y="19"
				width="189"
				height="15.8537"
				rx="7.92683"
				fill={color}
			/>
			<rect
				x="179"
				width="114"
				height="10.5691"
				rx="5.28455"
				fill={color}
			/>
		</svg>
	)
}

export default CustomIcon
