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
    public partial class subImage : Form
    {
        public subImage()
        {
            InitializeComponent();
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
            m_Picturebox_Canvas.Focus();
            if (m_Picturebox_Canvas.Focused == true)
            {
                if (mea.Delta > 0)
                {
                    ZoomInScroll();
                }
                else if (mea.Delta < 0)
                {
                    ZoomOutScroll();
                }
            }
        }
    }
}
