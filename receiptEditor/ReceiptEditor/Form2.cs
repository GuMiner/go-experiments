using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows.Forms;

namespace ReceiptEditor
{
    public partial class Form2 : Form
    {
        private SubImage subImage;

        public Form2()
        {
            InitializeComponent();
        }

        public Form2(SubImage subImage)
            : this()
        {
            this.subImage = subImage;
        }

        private void pictureBox1_MouseMove(object sender, MouseEventArgs e)
        {

        }

        private void pictureBox1_MouseDown(object sender, MouseEventArgs e)
        {

        }

        private void pictureBox1_MouseUp(object sender, MouseEventArgs e)
        {

        }

        private void subImage_Scroll(object sender, ScrollEventArgs e)
        {

        }

        protected override void OnMouseWheel(MouseEventArgs mea)
        {
            // subImageViewer.Canvas
            // m_Picturebox_Canvas.Focus();
            // if (m_Picturebox_Canvas.Focused == true)
            // {
            //     if (mea.Delta > 0)
            //     {
            //         ZoomInScroll();
            //     }
            //     else if (mea.Delta < 0)
            //     {
            //         ZoomOutScroll();
            //     }
            // }
        }

        private void imageBox_Paint(object sender, PaintEventArgs e)
        {
            e.Graphics.Clear(Color.LightSeaGreen);
            e.Graphics.DrawImage(subImage.Image,
                new Rectangle(0, 0, imageBox.Width, imageBox.Height),
                new Rectangle(subImage.MinPos, new Size(subImage.MaxPos.X - subImage.MinPos.X, subImage.MaxPos.Y - subImage.MinPos.Y)),
                GraphicsUnit.Pixel);
        }

        private void Form2_Resize(object sender, EventArgs e)
        {
            imageBox.Invalidate();
        }
    }
}
